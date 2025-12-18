package reward

import (
	"time"
)

type RewardEvent struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	StockID     int       `json:"stock_id"`
	Quantity    float64   `json:"quantity"`
	StockPrice  float64   `json:"stock_price"`
	TotalValue  float64   `json:"total_value"`
	EventType   string    `json:"event_type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRewardRequest struct {
	UserID         int     `json:"user_id" binding:"required"`
	StockSymbol    string  `json:"stock_symbol" binding:"required"`
	Quantity       float64 `json:"quantity" binding:"required,gt=0"`
	Description    string  `json:"description"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type RewardEventWithDetails struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserEmail   string    `json:"user_email"`
	StockID     int       `json:"stock_id"`
	StockSymbol string    `json:"stock_symbol"`
	StockName   string    `json:"stock_name"`
	Quantity    float64   `json:"quantity"`
	StockPrice  float64   `json:"stock_price"`
	TotalValue  float64   `json:"total_value"`
	EventType   string    `json:"event_type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type FeeConfiguration struct {
	ID          int       `json:"id"`
	FeeType     string    `json:"fee_type"`
	Percentage  float64   `json:"percentage"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaginatedRewardsResponse struct {
	Data       []RewardEventWithDetails `json:"data"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalCount int                      `json:"total_count"`
	TotalPages int                      `json:"total_pages"`
}
