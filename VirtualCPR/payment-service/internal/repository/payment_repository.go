package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/internal/models"
)

type PaymentRepository interface {
	CreatePayment(payment models.Payment) (models.Payment, error)
	GetPaymentByID(id string) (models.Payment, error)
	ListPaymentsByUserID(userID string) ([]models.Payment, error)
	UpdatePayment(payment models.Payment) error
}

type PostgresPaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{
		db: db,
	}
}

func (r *PostgresPaymentRepository) CreatePayment(payment models.Payment) (models.Payment, error) {
	query := `INSERT INTO payments (id, user_id, amount, currency, description, status, stripe_charge_id, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
              RETURNING id, user_id, amount, currency, description, status, stripe_charge_id, created_at, updated_at`
	if payment.ID == "" {
		payment.ID = uuid.New().String()
	}
	now := time.Now()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		payment.ID,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.Desc,
		payment.Status,
		payment.StripeChargeID,
		payment.CreatedAt,
		payment.UpdatedAt,
	).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.Desc,
		&payment.Status,
		&payment.StripeChargeID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return models.Payment{}, err
	}
	return payment, nil
}

func (r *PostgresPaymentRepository) GetPaymentByID(id string) (models.Payment, error) {
	query := `SELECT id, user_id, amount, currency, description, status, stripe_charge_id, created_at, updated_at
			  FROM payments
			  WHERE id = $1`
	var payment models.Payment
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.Desc,
		&payment.Status,
		&payment.StripeChargeID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return models.Payment{}, err
	}
	return payment, nil
}

func (r *PostgresPaymentRepository) ListPaymentsByUserID(userID string) ([]models.Payment, error) {
	query := `SELECT id, user_id, amount, currency, description, status, stripe_charge_id, created_at, updated_at
			  FROM payments
			  WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&payment.Amount,
			&payment.Currency,
			&payment.Desc,
			&payment.Status,
			&payment.StripeChargeID,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PostgresPaymentRepository) UpdatePayment(payment models.Payment) error {
	query := `UPDATE payments SET status = $1, stripe_charge_id = $2, updated_at = $3 WHERE id = $4`

	payment.UpdatedAt = time.Now()
	result, err := r.db.Exec(query, payment.Status, payment.StripeChargeID, payment.UpdatedAt, payment.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("payment not found")
	}
	return nil
}
