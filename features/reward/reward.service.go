package reward

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type RewardService struct {
	db *sql.DB
}

func NewRewardService(db *sql.DB) *RewardService {
	return &RewardService{db: db}
}

func (s *RewardService) CreateReward(req CreateRewardRequest) (*RewardEvent, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.Errorf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	var stockID int
	var stockPrice float64
	err = tx.QueryRow(`SELECT id, current_price FROM stocks WHERE symbol = $1 AND is_active = true`, 
		req.StockSymbol).Scan(&stockID, &stockPrice)
	if err != nil {
		logrus.Errorf("Failed to get stock details: %v", err)
		return nil, fmt.Errorf("stock not found or inactive")
	}

	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND is_active = true)`, 
		req.UserID).Scan(&userExists)
	if err != nil || !userExists {
		return nil, fmt.Errorf("user not found or inactive")
	}

	var duplicateExists bool
	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM reward_events 
			WHERE user_id = $1 
			AND stock_id = $2 
			AND quantity = $3 
			AND created_at > NOW() - INTERVAL '5 minutes'
		)
	`, req.UserID, stockID, req.Quantity).Scan(&duplicateExists)
	if err != nil {
		logrus.Errorf("Failed to check for duplicate reward: %v", err)
		return nil, err
	}
	if duplicateExists {
		return nil, fmt.Errorf("duplicate reward detected: similar reward was created within the last 5 minutes")
	}

	if req.IdempotencyKey != "" {
		var existingRewardID int
		err = tx.QueryRow(`
			SELECT id FROM reward_events 
			WHERE description LIKE $1
			AND created_at > NOW() - INTERVAL '1 hour'
			LIMIT 1
		`, "%idempotency:"+req.IdempotencyKey+"%").Scan(&existingRewardID)
		if err == nil {
			return nil, fmt.Errorf("duplicate request: idempotency key already used")
		} else if err != sql.ErrNoRows {
			logrus.Errorf("Failed to check idempotency key: %v", err)
			return nil, err
		}
	}

	totalValue := req.Quantity * stockPrice

	fees, err := s.getFeeConfigurations(tx)
	if err != nil {
		return nil, err
	}

	brokerageFee := totalValue * fees["BROKERAGE"]
	sttFee := totalValue * fees["STT"]
	gstFee := brokerageFee * fees["GST"]

	// Append idempotency key to description for tracking
	description := req.Description
	if req.IdempotencyKey != "" {
		description = fmt.Sprintf("%s [idempotency:%s]", req.Description, req.IdempotencyKey)
	}

	var rewardEvent RewardEvent
	err = tx.QueryRow(`
		INSERT INTO reward_events (user_id, stock_id, quantity, stock_price, total_value, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, stock_id, quantity, stock_price, total_value, event_type, status, description, created_at, updated_at
	`, req.UserID, stockID, req.Quantity, stockPrice, totalValue, description).Scan(
		&rewardEvent.ID, &rewardEvent.UserID, &rewardEvent.StockID, &rewardEvent.Quantity,
		&rewardEvent.StockPrice, &rewardEvent.TotalValue, &rewardEvent.EventType,
		&rewardEvent.Status, &rewardEvent.Description, &rewardEvent.CreatedAt, &rewardEvent.UpdatedAt,
	)
	if err != nil {
		logrus.Errorf("Failed to create reward event: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO ledger_entries (reward_event_id, user_id, entry_type, account_type, stock_id, quantity, description)
		VALUES ($1, $2, 'DEBIT', 'STOCK_UNITS', $3, $4, 'Stock reward credited')
	`, rewardEvent.ID, req.UserID, stockID, req.Quantity)
	if err != nil {
		logrus.Errorf("Failed to create stock ledger entry: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO ledger_entries (reward_event_id, user_id, entry_type, account_type, amount, description)
		VALUES ($1, $2, 'CREDIT', 'INR_CASH', $3, 'Cash outflow for stock purchase')
	`, rewardEvent.ID, req.UserID, totalValue)
	if err != nil {
		logrus.Errorf("Failed to create cash ledger entry: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO ledger_entries (reward_event_id, user_id, entry_type, account_type, amount, description)
		VALUES ($1, $2, 'CREDIT', 'BROKERAGE_FEE', $3, 'Brokerage fee')
	`, rewardEvent.ID, req.UserID, brokerageFee)
	if err != nil {
		logrus.Errorf("Failed to create brokerage fee ledger entry: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO ledger_entries (reward_event_id, user_id, entry_type, account_type, amount, description)
		VALUES ($1, $2, 'CREDIT', 'STT_FEE', $3, 'Securities Transaction Tax')
	`, rewardEvent.ID, req.UserID, sttFee)
	if err != nil {
		logrus.Errorf("Failed to create STT fee ledger entry: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO ledger_entries (reward_event_id, user_id, entry_type, account_type, amount, description)
		VALUES ($1, $2, 'CREDIT', 'GST_FEE', $3, 'GST on brokerage')
	`, rewardEvent.ID, req.UserID, gstFee)
	if err != nil {
		logrus.Errorf("Failed to create GST fee ledger entry: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO user_stock_holdings (user_id, stock_id, total_quantity, average_price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, stock_id) 
		DO UPDATE SET 
			total_quantity = user_stock_holdings.total_quantity + EXCLUDED.total_quantity,
			average_price = ((user_stock_holdings.total_quantity * user_stock_holdings.average_price) + 
							(EXCLUDED.total_quantity * EXCLUDED.average_price)) / 
							(user_stock_holdings.total_quantity + EXCLUDED.total_quantity),
			updated_at = CURRENT_TIMESTAMP
	`, req.UserID, stockID, req.Quantity, stockPrice)
	if err != nil {
		logrus.Errorf("Failed to update user stock holdings: %v", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		logrus.Errorf("Failed to commit transaction: %v", err)
		return nil, err
	}

	logrus.Infof("Reward created successfully: User %d received %.6f units of stock %d", 
		req.UserID, req.Quantity, stockID)
	return &rewardEvent, nil
}

func (s *RewardService) GetAllRewards(page, pageSize int) (*PaginatedRewardsResponse, error) {
	var totalCount int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM reward_events`).Scan(&totalCount)
	if err != nil {
		logrus.Errorf("Failed to count rewards: %v", err)
		return nil, err
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT 
			re.id, re.user_id, u.name as user_name, u.email as user_email,
			re.stock_id, s.symbol as stock_symbol, s.name as stock_name,
			re.quantity, re.stock_price, re.total_value,
			re.event_type, re.status, re.description, re.created_at
		FROM reward_events re
		JOIN users u ON re.user_id = u.id
		JOIN stocks s ON re.stock_id = s.id
		ORDER BY re.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		logrus.Errorf("Failed to query rewards: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rewards []RewardEventWithDetails
	for rows.Next() {
		var reward RewardEventWithDetails
		err := rows.Scan(
			&reward.ID, &reward.UserID, &reward.UserName, &reward.UserEmail,
			&reward.StockID, &reward.StockSymbol, &reward.StockName,
			&reward.Quantity, &reward.StockPrice, &reward.TotalValue,
			&reward.EventType, &reward.Status, &reward.Description, &reward.CreatedAt,
		)
		if err != nil {
			logrus.Errorf("Failed to scan reward: %v", err)
			return nil, err
		}
		rewards = append(rewards, reward)
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	return &PaginatedRewardsResponse{
		Data:       rewards,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (s *RewardService) GetRewardsByUserID(userID, page, pageSize int) (*PaginatedRewardsResponse, error) {
	var totalCount int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM reward_events WHERE user_id = $1`, userID).Scan(&totalCount)
	if err != nil {
		logrus.Errorf("Failed to count user rewards: %v", err)
		return nil, err
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT 
			re.id, re.user_id, u.name as user_name, u.email as user_email,
			re.stock_id, s.symbol as stock_symbol, s.name as stock_name,
			re.quantity, re.stock_price, re.total_value,
			re.event_type, re.status, re.description, re.created_at
		FROM reward_events re
		JOIN users u ON re.user_id = u.id
		JOIN stocks s ON re.stock_id = s.id
		WHERE re.user_id = $1
		ORDER BY re.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, userID, pageSize, offset)
	if err != nil {
		logrus.Errorf("Failed to query rewards: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rewards []RewardEventWithDetails
	for rows.Next() {
		var reward RewardEventWithDetails
		err := rows.Scan(
			&reward.ID, &reward.UserID, &reward.UserName, &reward.UserEmail,
			&reward.StockID, &reward.StockSymbol, &reward.StockName,
			&reward.Quantity, &reward.StockPrice, &reward.TotalValue,
			&reward.EventType, &reward.Status, &reward.Description, &reward.CreatedAt,
		)
		if err != nil {
			logrus.Errorf("Failed to scan reward: %v", err)
			return nil, err
		}
		rewards = append(rewards, reward)
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	return &PaginatedRewardsResponse{
		Data:       rewards,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (s *RewardService) AdjustReward(req AdjustRewardRequest) (*RewardEvent, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.Errorf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	var originalReward RewardEvent
	err = tx.QueryRow(`
		SELECT id, user_id, stock_id, quantity, stock_price, total_value, event_type, status
		FROM reward_events
		WHERE id = $1
	`, req.RewardEventID).Scan(
		&originalReward.ID, &originalReward.UserID, &originalReward.StockID,
		&originalReward.Quantity, &originalReward.StockPrice, &originalReward.TotalValue,
		&originalReward.EventType, &originalReward.Status,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("reward event not found")
	}
	if err != nil {
		logrus.Errorf("Failed to fetch reward event: %v", err)
		return nil, err
	}

	if originalReward.EventType == "ADJUSTMENT" {
		return nil, fmt.Errorf("cannot adjust an adjustment entry")
	}

	if req.AdjustmentType == "REFUND" && req.Quantity != originalReward.Quantity {
		return nil, fmt.Errorf("full refund must match original quantity: %.6f", originalReward.Quantity)
	}
	if req.AdjustmentType == "PARTIAL_REFUND" && req.Quantity >= originalReward.Quantity {
		return nil, fmt.Errorf("partial refund quantity must be less than original: %.6f", originalReward.Quantity)
	}

	var currentHoldings float64
	err = tx.QueryRow(`
		SELECT total_quantity FROM user_stock_holdings
		WHERE user_id = $1 AND stock_id = $2
	`, originalReward.UserID, originalReward.StockID).Scan(&currentHoldings)
	if err != nil {
		return nil, fmt.Errorf("user stock holdings not found")
	}
	if currentHoldings < req.Quantity {
		return nil, fmt.Errorf("insufficient holdings: user has %.6f, adjustment requires %.6f", currentHoldings, req.Quantity)
	}

	adjustmentValue := req.Quantity * originalReward.StockPrice
	description := fmt.Sprintf("%s for reward #%d: %s", req.AdjustmentType, req.RewardEventID, req.Reason)

	var adjustmentEvent RewardEvent
	err = tx.QueryRow(`
		INSERT INTO reward_events (user_id, stock_id, quantity, stock_price, total_value, event_type, description)
		VALUES ($1, $2, $3, $4, $5, 'ADJUSTMENT', $6)
		RETURNING id, user_id, stock_id, quantity, stock_price, total_value, event_type, status, description, created_at, updated_at
	`, originalReward.UserID, originalReward.StockID, req.Quantity, originalReward.StockPrice, adjustmentValue, description).Scan(
		&adjustmentEvent.ID, &adjustmentEvent.UserID, &adjustmentEvent.StockID,
		&adjustmentEvent.Quantity, &adjustmentEvent.StockPrice, &adjustmentEvent.TotalValue,
		&adjustmentEvent.EventType, &adjustmentEvent.Status, &adjustmentEvent.Description,
		&adjustmentEvent.CreatedAt, &adjustmentEvent.UpdatedAt,
	)
	if err != nil {
		logrus.Errorf("Failed to create adjustment event: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO ledger_entries (reward_event_id, user_id, entry_type, account_type, stock_id, quantity, description)
		VALUES ($1, $2, 'CREDIT', 'STOCK_UNITS', $3, $4, $5)
	`, adjustmentEvent.ID, originalReward.UserID, originalReward.StockID, -req.Quantity, description)
	if err != nil {
		logrus.Errorf("Failed to create stock ledger entry: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO ledger_entries (reward_event_id, user_id, entry_type, account_type, amount, description)
		VALUES ($1, $2, 'DEBIT', 'INR_CASH', $3, $4)
	`, adjustmentEvent.ID, originalReward.UserID, adjustmentValue, description)
	if err != nil {
		logrus.Errorf("Failed to create credit ledger entry: %v", err)
		return nil, err
	}

	_, err = tx.Exec(`
		UPDATE user_stock_holdings
		SET total_quantity = total_quantity - $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2 AND stock_id = $3
	`, req.Quantity, originalReward.UserID, originalReward.StockID)
	if err != nil {
		logrus.Errorf("Failed to update user stock holdings: %v", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		logrus.Errorf("Failed to commit adjustment transaction: %v", err)
		return nil, err
	}

	logrus.Infof("Reward adjusted successfully: User %d, %.6f units refunded for reward #%d",
		originalReward.UserID, req.Quantity, req.RewardEventID)
	return &adjustmentEvent, nil
}

func (s *RewardService) getFeeConfigurations(tx *sql.Tx) (map[string]float64, error) {
	rows, err := tx.Query(`SELECT fee_type, percentage FROM fee_configurations WHERE is_active = true`)
	if err != nil {
		logrus.Errorf("Failed to query fee configurations: %v", err)
		return nil, err
	}
	defer rows.Close()

	fees := make(map[string]float64)
	for rows.Next() {
		var feeType string
		var percentage float64
		if err := rows.Scan(&feeType, &percentage); err != nil {
			return nil, err
		}
		fees[feeType] = percentage
	}

	return fees, nil
}
