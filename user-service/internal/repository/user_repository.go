package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/user-service/internal/models"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByID(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	ListUsers() ([]models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(id string) error
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) CreateUser(user models.User) (models.User, error) {
	query := `INSERT INTO users (id, name, email, address, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, email, address, password,created_at, updated_at`

	// Generate UUID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRow(query, user.ID, user.Name, user.Email, user.Address, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserByID(id string) (models.User, error) {
	query := `SELECT id, name, email, address, password, created_at, updated_at FROM users WHERE id = $1`
	var user models.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (models.User, error) {
	query := `SELECT id, name, email, address, password, created_at, updated_at FROM users WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) ListUsers() ([]models.User, error) {
	query := `SELECT id, name, email, address, password, created_at, updated_at FROM users`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) UpdateUser(user models.User) error {
	query := `UPDATE users SET name = $1, email = $2, address = $3, password = $4, updated_at = $5 WHERE id = $6`

	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(query, user.Name, user.Email, user.Address, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *PostgresUserRepository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
