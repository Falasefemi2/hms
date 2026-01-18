package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(databaseURL string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	return &DB{pool: pool}, nil
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) HealthCheck(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) InitializeSchema(ctx context.Context) error {
	log.Println("Reading schema.sql file...")

	schemaSQL, err := readSchemaFile()
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	statements := parseSQLStatements(schemaSQL)

	for i, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}

		_, err := db.pool.Exec(ctx, stmt)
		if err != nil {
			return fmt.Errorf(
				"failed to execute statement %d: %w\nStatement: %s",
				i+1,
				err,
				stmt,
			)
		}
	}

	log.Println("Database schema initialized successfully")
	return nil
}

func readSchemaFile() (string, error) {
	// Try multiple possible locations for schema.sql
	possiblePaths := []string{
		"schema.sql",
		"./schema.sql",
		"../schema.sql",
		"../../schema.sql",
	}

	var schemaContent string
	var lastErr error

	for _, path := range possiblePaths {
		content, err := os.ReadFile(path)
		if err == nil {
			schemaContent = string(content)
			log.Printf("Loaded schema.sql from: %s\n", path)
			return schemaContent, nil
		}
		lastErr = err
	}

	return "", fmt.Errorf("schema.sql not found in expected locations: %w", lastErr)
}

func parseSQLStatements(sqlContent string) []string {
	var statements []string
	var currentStatement strings.Builder

	lines := strings.Split(sqlContent, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "--") {
			continue
		}

		// Add line to current statement
		currentStatement.WriteString(trimmedLine)
		currentStatement.WriteString("\n")

		// Check if this line ends a statement (ends with semicolon)
		if strings.HasSuffix(trimmedLine, ";") {
			stmt := currentStatement.String()
			if strings.TrimSpace(stmt) != "" {
				statements = append(statements, stmt)
			}
			currentStatement.Reset()
		}
	}

	// Add any remaining statement
	if currentStatement.Len() > 0 {
		stmt := currentStatement.String()
		if strings.TrimSpace(stmt) != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}
