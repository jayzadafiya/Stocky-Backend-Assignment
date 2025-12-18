package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"stocky-backend/config"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"is_active"`
}

type Stock struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Exchange     string `json:"exchange"`
	CurrentPrice string `json:"current_price"`
	IsActive     bool   `json:"is_active"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	config.InitLogger()

	logrus.Info("Starting database migration...")

	db, err := config.ConnectDatabase()
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase()

	if err := runMigrations(db); err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}

	if err := seedData(db); err != nil {
		logrus.Fatalf("Data seeding failed: %v", err)
	}

	logrus.Info("Migration and seeding completed successfully!")
}

func runMigrations(db *sql.DB) error {
	logrus.Info("Running SQL migrations...")

	migrationsPath := "migrations"
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		logrus.Infof("Running migration: %s", filename)

		content, err := ioutil.ReadFile(filepath.Join(migrationsPath, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		logrus.Infof("Migration %s completed successfully", filename)
	}

	return nil
}

func seedData(db *sql.DB) error {
	logrus.Info("Seeding data from JSON files...")

	if err := seedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := seedStocks(db); err != nil {
		return fmt.Errorf("failed to seed stocks: %w", err)
	}

	return nil
}

func seedUsers(db *sql.DB) error {
	logrus.Info("Seeding users...")

	content, err := ioutil.ReadFile("data/users.json")
	if err != nil {
		return fmt.Errorf("failed to read users.json: %w", err)
	}

	var users []User
	if err := json.Unmarshal(content, &users); err != nil {
		return fmt.Errorf("failed to parse users.json: %w", err)
	}

	for _, user := range users {
		_, err := db.Exec(`
			INSERT INTO users (email, name, phone, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (email) DO UPDATE SET
				name = EXCLUDED.name,
				phone = EXCLUDED.phone,
				is_active = EXCLUDED.is_active,
				updated_at = CURRENT_TIMESTAMP
		`, user.Email, user.Name, user.Phone, user.IsActive)

		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.Email, err)
		}
	}

	logrus.Infof("Successfully seeded %d users", len(users))
	return nil
}

func seedStocks(db *sql.DB) error {
	logrus.Info("Seeding stocks...")

	content, err := ioutil.ReadFile("data/stocks.json")
	if err != nil {
		return fmt.Errorf("failed to read stocks.json: %w", err)
	}

	var stocks []Stock
	if err := json.Unmarshal(content, &stocks); err != nil {
		return fmt.Errorf("failed to parse stocks.json: %w", err)
	}

	for _, stock := range stocks {
		price, err := strconv.ParseFloat(stock.CurrentPrice, 64)
		if err != nil {
			logrus.Warnf("Invalid price for stock %s: %v", stock.Symbol, err)
			continue
		}

		_, err = db.Exec(`
			INSERT INTO stocks (symbol, name, exchange, current_price, is_active)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (symbol) DO UPDATE SET
				name = EXCLUDED.name,
				exchange = EXCLUDED.exchange,
				current_price = EXCLUDED.current_price,
				is_active = EXCLUDED.is_active,
				updated_at = CURRENT_TIMESTAMP
		`, stock.Symbol, stock.Name, stock.Exchange, price, stock.IsActive)

		if err != nil {
			return fmt.Errorf("failed to insert stock %s: %w", stock.Symbol, err)
		}
	}

	logrus.Infof("Successfully seeded %d stocks", len(stocks))
	return nil
}
