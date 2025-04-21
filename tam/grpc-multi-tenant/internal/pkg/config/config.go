package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg/logger"
)

// Config holds all configuration for the application
type Config struct {
	// Database configuration
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBDSN      string

	// Server configuration
	ServerPort string

	// Logging configuration
	LogLevel logger.LogLevel
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Try to load .env file if it exists
	// Ignore errors since .env file is optional
	// This is useful for local development
	_, err := os.Stat(".env")
	if err == nil {
		godotenv.Load(".env")
	} else {
		// Try to look in the root of the project
		workDir, err := os.Getwd()
		if err == nil {
			envPath := filepath.Join(workDir, ".env")
			godotenv.Load(envPath)
		}
	}

	cfg := &Config{
		// Database config with defaults
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "database"),

		// Server config with defaults
		ServerPort: getEnv("SERVER_PORT", "50051"),

		// Logging config with defaults
		LogLevel: logger.LogLevel(getEnv("LOG_LEVEL", string(logger.Info))),
	}

	// Construct the database connection string
	cfg.DBDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	return cfg, nil
}

// getEnv retrieves an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
