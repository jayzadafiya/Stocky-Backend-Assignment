package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "assignment"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func ConnectDatabase() (*sql.DB, error) {
	config := LoadDatabaseConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.Errorf("Failed to open database connection: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logrus.Errorf("Failed to ping database: %v", err)
		return nil, err
	}

	logrus.Info("Database connection established successfully")
	DB = db
	return db, nil
}

func CloseDatabase() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			logrus.Errorf("Error closing database: %v", err)
		} else {
			logrus.Info("Database connection closed")
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
