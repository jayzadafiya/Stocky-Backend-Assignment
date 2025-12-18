package user

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type TodayStockReward struct {
	ID          int       `json:"id"`
	StockSymbol string    `json:"stock_symbol"`
	StockName   string    `json:"stock_name"`
	Quantity    float64   `json:"quantity"`
	StockPrice  float64   `json:"stock_price"`
	TotalValue  float64   `json:"total_value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginatedUsersResponse struct {
	Data       []User `json:"data"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalCount int    `json:"total_count"`
	TotalPages int    `json:"total_pages"`
}

type PaginatedStockRewardsResponse struct {
	Data       []TodayStockReward `json:"data"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalCount int                `json:"total_count"`
	TotalPages int                `json:"total_pages"`
}

type HistoricalINRValue struct {
	Date       string  `json:"date"`
	TotalValue float64 `json:"total_value"`
	RewardCount int    `json:"reward_count"`
}

type StockRewardSummary struct {
	StockSymbol string  `json:"stock_symbol"`
	StockName   string  `json:"stock_name"`
	TotalShares float64 `json:"total_shares"`
}

type UserStats struct {
	TodayRewards            []StockRewardSummary `json:"today_rewards"`
	CurrentPortfolioValue   float64              `json:"current_portfolio_value"`
}

type PortfolioHolding struct {
	StockSymbol    string  `json:"stock_symbol"`
	StockName      string  `json:"stock_name"`
	TotalQuantity  float64 `json:"total_quantity"`
	AveragePrice   float64 `json:"average_price"`
	CurrentPrice   float64 `json:"current_price"`
	CurrentValue   float64 `json:"current_value"`
	InvestmentCost float64 `json:"investment_cost"`
	ProfitLoss     float64 `json:"profit_loss"`
}

type PaginatedHistoricalINRResponse struct {
	Data       []HistoricalINRValue `json:"data"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalCount int                  `json:"total_count"`
	TotalPages int                  `json:"total_pages"`
}

type PaginatedPortfolioResponse struct {
	Data             []PortfolioHolding `json:"data"`
	Page             int                `json:"page"`
	PageSize         int                `json:"page_size"`
	TotalCount       int                `json:"total_count"`
	TotalPages       int                `json:"total_pages"`
	TotalPortfolioValue float64         `json:"total_portfolio_value"`
}
