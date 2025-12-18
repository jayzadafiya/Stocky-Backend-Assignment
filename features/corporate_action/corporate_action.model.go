package corporate_action

import (
	"time"
)

type CorporateActionType string

const (
	ActionStockSplit CorporateActionType = "STOCK_SPLIT"
	ActionMerger     CorporateActionType = "MERGER"
	ActionDelisting  CorporateActionType = "DELISTING"
)

type CorporateAction struct {
	ID             int                 `json:"id"`
	StockID        int                 `json:"stock_id"`
	StockSymbol    string              `json:"stock_symbol"`
	ActionType     CorporateActionType `json:"action_type"`
	SplitRatio     float64             `json:"split_ratio,omitempty"`
	MergerToStockID int                `json:"merger_to_stock_id,omitempty"`
	MergerRatio    float64             `json:"merger_ratio,omitempty"`
	EffectiveDate  time.Time           `json:"effective_date"`
	Status         string              `json:"status"`
	Description    string              `json:"description"`
	CreatedAt      time.Time           `json:"created_at"`
	ProcessedAt    *time.Time          `json:"processed_at,omitempty"`
}

type CreateCorporateActionRequest struct {
	StockSymbol     string              `json:"stock_symbol" binding:"required"`
	ActionType      CorporateActionType `json:"action_type" binding:"required"`
	SplitRatio      float64             `json:"split_ratio,omitempty"`
	MergerToSymbol  string              `json:"merger_to_symbol,omitempty"`
	MergerRatio     float64             `json:"merger_ratio,omitempty"`
	EffectiveDate   string              `json:"effective_date" binding:"required"`
	Description     string              `json:"description"`
}

type CorporateActionResponse struct {
	ID                int                 `json:"id"`
	StockSymbol       string              `json:"stock_symbol"`
	ActionType        CorporateActionType `json:"action_type"`
	SplitRatio        float64             `json:"split_ratio,omitempty"`
	MergerToSymbol    string              `json:"merger_to_symbol,omitempty"`
	MergerRatio       float64             `json:"merger_ratio,omitempty"`
	EffectiveDate     string              `json:"effective_date"`
	Status            string              `json:"status"`
	Description       string              `json:"description"`
	AffectedUsers     int                 `json:"affected_users"`
	CreatedAt         time.Time           `json:"created_at"`
	ProcessedAt       *time.Time          `json:"processed_at,omitempty"`
}

type PaginatedCorporateActionsResponse struct {
	Data       []CorporateActionResponse `json:"data"`
	Page       int                       `json:"page"`
	PageSize   int                       `json:"page_size"`
	TotalCount int                       `json:"total_count"`
	TotalPages int                       `json:"total_pages"`
}
