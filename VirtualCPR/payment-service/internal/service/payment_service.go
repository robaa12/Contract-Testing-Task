package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/internal/models"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/internal/repository"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) CreatePayment(request models.CreatePaymentRequest) (models.PaymentResponse, error) {
	payment := models.Payment{
		ID:        uuid.New().String(),
		UserID:    request.UserID,
		Amount:    request.Amount,
		Currency:  request.Currency,
		Desc:      request.Desc,
		Status:    models.PaymentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdPayment, err := s.repo.CreatePayment(payment)
	if err != nil {
		return models.PaymentResponse{}, err
	}

	// Create a Stripe charge
	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(request.Amount),
		Currency:    stripe.String(request.Currency),
		Description: stripe.String(request.Desc),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Confirm: stripe.Bool(true),
		Metadata: map[string]string{
			"payment_id": createdPayment.ID,
			"user_id":    request.UserID,
		},
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		createdPayment.Status = models.PaymentStatusFailed
		_ = s.repo.UpdatePayment(createdPayment)
		return models.PaymentResponse{}, err
	}

	createdPayment.Status = models.PaymentStatusSucceeded
	createdPayment.StripeChargeID = pi.ID // Store PaymentIntent ID instead of Charge ID
	err = s.repo.UpdatePayment(createdPayment)
	if err != nil {
		return models.PaymentResponse{}, err
	}

	updatedPayment, err := s.repo.GetPaymentByID(createdPayment.ID)
	if err != nil {
		return models.PaymentResponse{}, err
	}
	return updatedPayment.ToPaymentResponse(), nil
}

func (s *PaymentService) GetPaymentByID(id string) (models.PaymentResponse, error) {
	payment, err := s.repo.GetPaymentByID(id)
	if err != nil {
		return models.PaymentResponse{}, err
	}
	return payment.ToPaymentResponse(), nil
}

func (s *PaymentService) ListPaymentsByUserID(userID string) ([]models.PaymentResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	payments, err := s.repo.ListPaymentsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var paymentResponses []models.PaymentResponse
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, payment.ToPaymentResponse())
	}
	return paymentResponses, nil
}
