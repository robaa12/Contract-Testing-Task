package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/internal/models"
)

type OrderRepository interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrderByID(orderID string) (models.Order, error)
	GetOrdersByUserID(userID string) ([]models.Order, error)
	ListOrders() ([]models.Order, error)
	UpdateOrderStatus(orderID, status string) (models.Order, error)
	DeleteOrder(orderID string) error
}

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) PostgresOrderRepository {
	return PostgresOrderRepository{
		db: db,
	}
}

func (r *PostgresOrderRepository) CreateOrder(order models.Order) (models.Order, error) {
	query := `INSERT INTO orders (id, user_id, products, total_amount, status, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING id, user_id, products, total_amount, status, created_at, updated_at`

	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	now := time.Now()
	order.CreatedAt = now
	order.UpdateAt = now

	if order.Status == "" {
		order.Status = "pending"
	}
	// Serialize products to JSON
	productsJSON, err := json.Marshal(order.Products)
	if err != nil {
		return models.Order{}, err
	}

	var returnedProductsJSON []byte
	err = r.db.QueryRow(query, order.ID, order.UserID, productsJSON, order.TotalAmount, order.Status, order.CreatedAt, order.UpdateAt).Scan(&order.ID, &order.UserID, &returnedProductsJSON, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdateAt)

	if err != nil {
		return models.Order{}, err
	}

	if err = json.Unmarshal(returnedProductsJSON, &order.Products); err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *PostgresOrderRepository) GetOrderByID(orderID string) (models.Order, error) {
	query := `SELECT id, user_id, products, total_amount, status, created_at, updated_at
			  FROM orders
			  WHERE id = $1`

	var order models.Order
	var productsJSON []byte
	err := r.db.QueryRow(query, orderID).Scan(&order.ID, &order.UserID, &productsJSON, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdateAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, errors.New("order not found")
		}
		return models.Order{}, err
	}

	if err = json.Unmarshal(productsJSON, &order.Products); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *PostgresOrderRepository) GetOrdersByUserID(userID string) ([]models.Order, error) {
	query := `SELECT id, user_id, products, total_amount, status, created_at, updated_at
			  FROM orders
			  WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var productsJSON []byte
		if err = rows.Scan(&order.ID, &order.UserID, &productsJSON, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdateAt); err != nil {
			return nil, err
		}

		if err = json.Unmarshal(productsJSON, &order.Products); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) ListOrders() ([]models.Order, error) {
	query := `SELECT id, user_id, products, total_amount, status, created_at, updated_at
			  FROM orders`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var productsJSON []byte
		if err = rows.Scan(&order.ID, &order.UserID, &productsJSON, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdateAt); err != nil {
			return nil, err
		}

		if err = json.Unmarshal(productsJSON, &order.Products); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PostgresOrderRepository) UpdateOrderStatus(id, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *PostgresOrderRepository) DeleteOrder(id string) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}
