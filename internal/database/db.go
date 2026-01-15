package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

type DB struct {
	conn *sql.DB
}

func NewDB(databaseURL string) (*DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection succssfully")

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *DB) HealthCheck() error {
	return db.conn.Ping()
}

func (db *DB) InitializeSchema() error {
	log.Println("Reading schema.sql file...")

	schemaSQL, err := readSchemaFile()
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	log.Println("Initializing database schema...")

	statements := parseSQLStatements(schemaSQL)

	log.Printf("Found %d SQL statements to execute\n", len(statements))
	for i, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue // Skip empty statements
		}

		log.Printf("Executing statement %d/%d...", i+1, len(statements))

		_, err := db.conn.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement %d: %w\nStatement: %s", i+1, err, stmt)
		}
	}

	log.Println("Database schema initialized successfully!")
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
