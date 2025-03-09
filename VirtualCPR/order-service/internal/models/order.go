package models

import "time"

type Order struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Products    []Product `json:"products"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type CreateOrderRequest struct {
	UserID   string    `json:"user_id" binding:"required"`
	Products []Product `json:"products" binding:"required,min=1"`
}

type UpdateOrderRequest struct {
	Status string `json:"status" binding:"required,oneof=pending cancelled completed"`
}

type OrderResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserEmail   string    `json:"user_email"`
	Products    []Product `json:"products"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

func CalculateTotalAmount(products []Product) float64 {
	var totalAmount float64
	for _, p := range products {
		totalAmount += p.Price * float64(p.Quantity)
	}
	return totalAmount
}
