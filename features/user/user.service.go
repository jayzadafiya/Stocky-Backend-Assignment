package user

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetAllUsers(page, pageSize int) (*PaginatedUsersResponse, error) {
	var totalCount int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalCount)
	if err != nil {
		logrus.Errorf("Failed to count users: %v", err)
		return nil, err
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, email, name, phone, is_active, created_at, updated_at 
			  FROM users ORDER BY created_at DESC
			  LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		logrus.Errorf("Failed to query users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Phone, 
						 &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			logrus.Errorf("Failed to scan user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	return &PaginatedUsersResponse{
		Data:       users,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserByID(id int) (*User, error) {
	query := `SELECT id, email, name, phone, is_active, created_at, updated_at 
			  FROM users WHERE id = $1`

	var user User
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Phone,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		logrus.Errorf("Failed to query user: %v", err)
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetTodayStockRewards(userID, page, pageSize int) (*PaginatedStockRewardsResponse, error) {
	var totalCount int
	err := s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM reward_events 
		WHERE user_id = $1 AND DATE(created_at) = CURRENT_DATE
	`, userID).Scan(&totalCount)
	if err != nil {
		logrus.Errorf("Failed to count today's stock rewards: %v", err)
		return nil, err
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT 
			re.id,
			s.symbol as stock_symbol,
			s.name as stock_name,
			re.quantity,
			re.stock_price,
			re.total_value,
			re.description,
			re.created_at
		FROM reward_events re
		JOIN stocks s ON re.stock_id = s.id
		WHERE re.user_id = $1 
		AND DATE(re.created_at) = CURRENT_DATE
		ORDER BY re.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, userID, pageSize, offset)
	if err != nil {
		logrus.Errorf("Failed to query today's stock rewards: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rewards []TodayStockReward
	for rows.Next() {
		var reward TodayStockReward
		err := rows.Scan(
			&reward.ID,
			&reward.StockSymbol,
			&reward.StockName,
			&reward.Quantity,
			&reward.StockPrice,
			&reward.TotalValue,
			&reward.Description,
			&reward.CreatedAt,
		)
		if err != nil {
			logrus.Errorf("Failed to scan stock reward: %v", err)
			return nil, err
		}
		rewards = append(rewards, reward)
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	return &PaginatedStockRewardsResponse{
		Data:       rewards,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}
