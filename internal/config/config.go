package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Database
	DatabaseURL string
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string

	// Server
	ServerPort  int
	ServerHost  string
	Environment string

	// JWT
	JWTSecret string
	JWTExpiry string

	// Logging
	LogLevel string
}

// LoadConfig loads configuration from environment variables
// It reads from .env file first (if it exists), then from environment variables
func LoadConfig() *Config {
	// Try to load .env file (but don't fail if it doesn't exist)
	// This allows using environment variables directly
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("  .env file not found - using environment variables instead")
		log.Println("    Make sure to set environment variables:")
		log.Println("    - DATABASE_URL (required)")
		log.Println("    - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME (if DATABASE_URL not set)")
		log.Println("    - Other variables (see .env.example)")
	} else {
		log.Println(".env file loaded successfully")
	}

	// Create config struct
	cfg := &Config{
		// Database configuration
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvAsInt("DB_PORT", 5432),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "hospitaldb"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "disable"),

		// Server configuration
		ServerPort:  getEnvAsInt("PORT", getEnvAsInt("SERVER_PORT", 8080)),
		ServerHost:  getEnv("SERVER_HOST", "0.0.0.0"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// JWT configuration
		JWTSecret: getEnv("JWT_SECRET", "default-secret-change-in-production"),
		JWTExpiry: getEnv("JWT_EXPIRY", "24h"),

		// Logging configuration
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return cfg
}

// GetDSN returns the PostgreSQL connection string
// Uses DATABASE_URL if available, otherwise builds from individual components
func (c *Config) GetDSN() string {
	// If DATABASE_URL is set, use it
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	// Otherwise build from components
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
	return dsn
}

// GetServerAddress returns the server address (host:port)
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort)
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	if c.DatabaseURL == "" && c.DBHost == "" {
		return fmt.Errorf("database configuration missing: either DATABASE_URL or DB_HOST is required")
	}

	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required but not set")
	}

	if c.JWTSecret == "default-secret-change-in-production" {
		log.Println("⚠️  WARNING: Using default JWT secret - set JWT_SECRET in .env for production!")
	}

	return nil
}

// Helper functions

// getEnv retrieves an environment variable with a default fallback
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvAsInt retrieves an environment variable as an integer with a default fallback
func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
