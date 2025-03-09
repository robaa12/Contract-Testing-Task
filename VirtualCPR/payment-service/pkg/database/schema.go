package database

import (
	"database/sql"
	"fmt"
)

func SetupSchema(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS payments (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			amount BIGINT NOT NULL,
			currency VARCHAR(3) NOT NULL,
			description TEXT,
			status VARCHAR(20) NOT NULL,
			stripe_charge_id VARCHAR(255),
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		`
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create payments table: %w", err)
		}

		return nil
}
