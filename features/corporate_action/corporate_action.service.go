package corporate_action

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type CorporateActionService struct {
	db *sql.DB
}

func NewCorporateActionService(db *sql.DB) *CorporateActionService {
	return &CorporateActionService{db: db}
}

func (s *CorporateActionService) CreateCorporateAction(req CreateCorporateActionRequest) (*CorporateActionResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.Errorf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	var stockID int
	err = tx.QueryRow(`SELECT id FROM stocks WHERE symbol = $1 AND is_active = true`, req.StockSymbol).Scan(&stockID)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %s", req.StockSymbol)
	}

	var mergerToStockID *int
	var mergerToSymbol string
	if req.ActionType == ActionMerger && req.MergerToSymbol != "" {
		var targetID int
		err = tx.QueryRow(`SELECT id FROM stocks WHERE symbol = $1 AND is_active = true`, req.MergerToSymbol).Scan(&targetID)
		if err != nil {
			return nil, fmt.Errorf("merger target stock not found: %s", req.MergerToSymbol)
		}
		mergerToStockID = &targetID
		mergerToSymbol = req.MergerToSymbol
	}

	var actionID int
	var createdAt time.Time
	err = tx.QueryRow(`
		INSERT INTO corporate_actions (stock_id, action_type, split_ratio, merger_to_stock_id, merger_ratio, effective_date, description, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING')
		RETURNING id, created_at
	`, stockID, req.ActionType, nullFloat64(req.SplitRatio), mergerToStockID, nullFloat64(req.MergerRatio), req.EffectiveDate, req.Description).Scan(&actionID, &createdAt)
	
	if err != nil {
		logrus.Errorf("Failed to create corporate action: %v", err)
		return nil, err
	}

	var affectedUsers int
	err = tx.QueryRow(`SELECT COUNT(DISTINCT user_id) FROM user_stock_holdings WHERE stock_id = $1 AND total_quantity > 0`, stockID).Scan(&affectedUsers)
	if err != nil {
		affectedUsers = 0
	}

	if err = tx.Commit(); err != nil {
		logrus.Errorf("Failed to commit transaction: %v", err)
		return nil, err
	}

	return &CorporateActionResponse{
		ID:             actionID,
		StockSymbol:    req.StockSymbol,
		ActionType:     req.ActionType,
		SplitRatio:     req.SplitRatio,
		MergerToSymbol: mergerToSymbol,
		MergerRatio:    req.MergerRatio,
		EffectiveDate:  req.EffectiveDate,
		Status:         "PENDING",
		Description:    req.Description,
		AffectedUsers:  affectedUsers,
		CreatedAt:      createdAt,
	}, nil
}

func (s *CorporateActionService) ProcessCorporateAction(actionID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var action CorporateAction
	var splitRatio, mergerRatio sql.NullFloat64
	var mergerToStockID sql.NullInt32
	
	var status string
	err = tx.QueryRow(`
		SELECT id, stock_id, action_type, split_ratio, merger_to_stock_id, merger_ratio, status
		FROM corporate_actions WHERE id = $1
	`, actionID).Scan(&action.ID, &action.StockID, &action.ActionType, &splitRatio, &mergerToStockID, &mergerRatio, &status)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("corporate action not found")
		}
		logrus.Errorf("Failed to fetch corporate action: %v", err)
		return fmt.Errorf("failed to fetch corporate action: %v", err)
	}
	
	if status == "COMPLETED" {
		return fmt.Errorf("corporate action already processed")
	}
	
	action.Status = status
	
	if splitRatio.Valid {
		action.SplitRatio = splitRatio.Float64
	}
	if mergerRatio.Valid {
		action.MergerRatio = mergerRatio.Float64
	}
	if mergerToStockID.Valid {
		action.MergerToStockID = int(mergerToStockID.Int32)
	}

	switch action.ActionType {
	case ActionStockSplit:
		err = s.processStockSplit(tx, action.StockID, action.SplitRatio)
	case ActionMerger:
		err = s.processMerger(tx, action.StockID, action.MergerToStockID, action.MergerRatio)
	case ActionDelisting:
		err = s.processDelisting(tx, action.StockID)
	default:
		return fmt.Errorf("unknown action type: %s", action.ActionType)
	}

	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE corporate_actions SET status = 'COMPLETED', processed_at = NOW() WHERE id = $1`, actionID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	logrus.Infof("Corporate action processed successfully: ID %d, Type %s", actionID, action.ActionType)
	return nil
}

func (s *CorporateActionService) processStockSplit(tx *sql.Tx, stockID int, splitRatio float64) error {
	_, err := tx.Exec(`
		UPDATE user_stock_holdings 
		SET total_quantity = total_quantity * $1,
		    average_price = average_price / $1,
		    updated_at = NOW()
		WHERE stock_id = $2 AND total_quantity > 0
	`, splitRatio, stockID)

	if err != nil {
		logrus.Errorf("Failed to process stock split: %v", err)
		return err
	}

	_, err = tx.Exec(`
		UPDATE stocks 
		SET current_price = current_price / $1,
		    updated_at = NOW()
		WHERE id = $2
	`, splitRatio, stockID)

	return err
}

func (s *CorporateActionService) processMerger(tx *sql.Tx, fromStockID int, toStockID int, mergerRatio float64) error {
	_, err := tx.Exec(`
		INSERT INTO user_stock_holdings (user_id, stock_id, total_quantity, average_price)
		SELECT user_id, $2, total_quantity * $3, average_price / $3
		FROM user_stock_holdings
		WHERE stock_id = $1 AND total_quantity > 0
		ON CONFLICT (user_id, stock_id) 
		DO UPDATE SET 
			total_quantity = user_stock_holdings.total_quantity + EXCLUDED.total_quantity,
			average_price = ((user_stock_holdings.total_quantity * user_stock_holdings.average_price) + 
			                (EXCLUDED.total_quantity * EXCLUDED.average_price)) / 
			                (user_stock_holdings.total_quantity + EXCLUDED.total_quantity),
			updated_at = NOW()
	`, fromStockID, toStockID, mergerRatio)

	if err != nil {
		logrus.Errorf("Failed to process merger: %v", err)
		return err
	}

	_, err = tx.Exec(`
		UPDATE user_stock_holdings 
		SET total_quantity = 0, updated_at = NOW()
		WHERE stock_id = $1
	`, fromStockID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE stocks SET is_active = false, updated_at = NOW() WHERE id = $1`, fromStockID)
	return err
}

func (s *CorporateActionService) processDelisting(tx *sql.Tx, stockID int) error {
	_, err := tx.Exec(`
		UPDATE user_stock_holdings 
		SET total_quantity = 0, updated_at = NOW()
		WHERE stock_id = $1 AND total_quantity > 0
	`, stockID)

	if err != nil {
		logrus.Errorf("Failed to process delisting: %v", err)
		return err
	}

	_, err = tx.Exec(`UPDATE stocks SET is_active = false, updated_at = NOW() WHERE id = $1`, stockID)
	return err
}

func (s *CorporateActionService) GetAllCorporateActions(page, pageSize int) (*PaginatedCorporateActionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var totalCount int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM corporate_actions`).Scan(&totalCount)
	if err != nil {
		logrus.Errorf("Failed to count corporate actions: %v", err)
		return nil, err
	}

	offset := (page - 1) * pageSize
	totalPages := (totalCount + pageSize - 1) / pageSize

	query := `
		SELECT 
			ca.id, s.symbol, ca.action_type, ca.split_ratio, ca.merger_ratio,
			COALESCE(s2.symbol, '') as merger_to_symbol,
			TO_CHAR(ca.effective_date, 'YYYY-MM-DD') as effective_date,
			ca.status, ca.description, ca.created_at, ca.processed_at,
			(SELECT COUNT(DISTINCT user_id) FROM user_stock_holdings WHERE stock_id = ca.stock_id AND total_quantity > 0) as affected_users
		FROM corporate_actions ca
		JOIN stocks s ON ca.stock_id = s.id
		LEFT JOIN stocks s2 ON ca.merger_to_stock_id = s2.id
		ORDER BY ca.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		logrus.Errorf("Failed to query corporate actions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var actions []CorporateActionResponse
	for rows.Next() {
		var action CorporateActionResponse
		var splitRatio, mergerRatio sql.NullFloat64
		err := rows.Scan(
			&action.ID, &action.StockSymbol, &action.ActionType,
			&splitRatio, &mergerRatio, &action.MergerToSymbol,
			&action.EffectiveDate, &action.Status, &action.Description,
			&action.CreatedAt, &action.ProcessedAt, &action.AffectedUsers,
		)
		if err != nil {
			logrus.Errorf("Failed to scan corporate action: %v", err)
			return nil, err
		}
		
		if splitRatio.Valid {
			action.SplitRatio = splitRatio.Float64
		}
		if mergerRatio.Valid {
			action.MergerRatio = mergerRatio.Float64
		}
		
		actions = append(actions, action)
	}

	return &PaginatedCorporateActionsResponse{
		Data:       actions,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func nullFloat64(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}
