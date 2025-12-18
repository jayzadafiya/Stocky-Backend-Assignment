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
	var stockName string
	err = tx.QueryRow(`SELECT id, current_price, name FROM stocks WHERE symbol = $1 AND is_active = true`, 
		req.StockSymbol).Scan(&stockID, &stockPrice, &stockName)
	if err != nil {
		logrus.Errorf("Failed to get stock details: %v", err)
		var isDelisted bool
		err2 := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM stocks WHERE symbol = $1 AND is_active = false)`, 
			req.StockSymbol).Scan(&isDelisted)
		if err2 == nil && isDelisted {
			return nil, fmt.Errorf("stock '%s' is delisted and cannot receive new rewards", req.StockSymbol)
		}
		return nil, fmt.Errorf("stock not found or inactive")
	}

	var pendingAction string
	err = tx.QueryRow(`
		SELECT action_type FROM corporate_actions 
		WHERE stock_id = $1 AND status = 'PENDING' 
		AND effective_date <= CURRENT_DATE
		LIMIT 1
	`, stockID).Scan(&pendingAction)
	if err == nil {
		return nil, fmt.Errorf("stock '%s' has a pending %s corporate action. Please process it before issuing new rewards", req.StockSymbol, pendingAction)
	} else if err != sql.ErrNoRows {
		logrus.Errorf("Failed to check corporate actions: %v", err)
		return nil, err
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
