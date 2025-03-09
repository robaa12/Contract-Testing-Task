package models

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID             string        `json:"id"`
	UserID         string        `json:"user_id"`
	Amount         int64         `json:"amount"`
	Currency       string        `json:"currency"`
	Desc           string        `json:"desc,omitempty"`
	Status         PaymentStatus `json:"status"`
	StripeChargeID string        `json:"stripe_charge_id,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type CreatePaymentRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required"`
	Currency  string `json:"currency" binding:"required"`
	Desc      string `json:"desc,omitempty"`
	CardToken string `json:"card_token" binding:"required"`
}

type PaymentResponse struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	Amount    int64         `json:"amount"`
	Currency  string        `json:"currency"`
	Desc      string        `json:"desc,omitempty"`
	Status    PaymentStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (p *Payment) ToPaymentResponse() PaymentResponse {
	return PaymentResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Desc:      p.Desc,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
