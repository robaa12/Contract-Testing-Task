package database

import (
	"database/sql"
	"fmt"
)

func SetupSchema(db *sql.DB) error {
	// Create users table
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		address TEXT,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}
