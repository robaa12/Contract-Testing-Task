package database

import (
	"database/sql"
	"fmt"
)

func SetupSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
	id VARCHAR(36) PRIMARY KEY,
	user_id VARCHAR(36) NOT NULL,
	products JSONB NOT NULL,
	total_amount DECIMAL(10,2) NOT NULL,
	status VARCHAR(20) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}
	return nil
}
