package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/internal/models"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/internal/repository"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/pkg/client"
)

type OrderService struct {
	repo       repository.PostgresOrderRepository
	userClient client.UserClient
}

func NewOrderService(repo repository.PostgresOrderRepository, userClient client.UserClient) *OrderService {
	return &OrderService{
		repo:       repo,
		userClient: userClient,
	}
}

func (s *OrderService) CreateOrder(req models.CreateOrderRequest) (models.OrderResponse, error) {
	user, err := s.userClient.ValidateUser(req.UserID)
	if err != nil {
		return models.OrderResponse{}, err
	}

	totalAmount := models.CalculateTotalAmount(req.Products)

	order := models.Order{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		Products:    req.Products,
		TotalAmount: totalAmount,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}

	createdOrder, err := s.repo.CreateOrder(order)
	if err != nil {
		return models.OrderResponse{}, err
	}
	return models.OrderResponse{
		ID:          createdOrder.ID,
		UserID:      createdOrder.UserID,
		UserName:    user.Name,
		UserEmail:   user.Email,
		Products:    createdOrder.Products,
		TotalAmount: createdOrder.TotalAmount,
		Status:      createdOrder.Status,
		CreatedAt:   createdOrder.CreatedAt,
		UpdateAt:    createdOrder.UpdateAt,
	}, nil
}

func (s *OrderService) GetOrder(orderID string) (models.OrderResponse, error) {
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return models.OrderResponse{}, err
	}

	user, err := s.userClient.ValidateUser(order.UserID)
	if err != nil {
		return models.OrderResponse{
			ID:          order.ID,
			Products:    order.Products,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdateAt:    order.UpdateAt,
		}, nil
	}

	return models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		UserName:    user.Name,
		UserEmail:   user.Email,
		Products:    order.Products,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdateAt:    order.UpdateAt,
	}, nil
}

func (s *OrderService) GetOrderByUserID(userID string) ([]models.OrderResponse, error) {
	user, err := s.userClient.ValidateUser(userID)
	if err != nil {
		return nil, err
	}
	orders, err := s.repo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	var orderResponses []models.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, models.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			UserName:    user.Name,
			UserEmail:   user.Email,
			Products:    order.Products,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdateAt:    order.UpdateAt,
		})
	}
	return orderResponses, nil
}

func (s *OrderService) ListOrders() ([]models.OrderResponse, error) {
	orders, err := s.repo.ListOrders()
	if err != nil {
		return nil, err
	}
	var orderResponses []models.OrderResponse
	for _, order := range orders {
		orderResponse := models.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			Products:    order.Products,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdateAt:    order.UpdateAt,
		}
		user, err := s.userClient.ValidateUser(order.UserID)
		if err == nil {
			orderResponse.UserName = user.Name
			orderResponse.UserEmail = user.Email
		}
		orderResponses = append(orderResponses, orderResponse)
	}
	return orderResponses, nil
}

func (s *OrderService) UpdateOrderStatus(id, status string) (models.OrderResponse, error) {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return models.OrderResponse{}, err
	}

	validStatuses := map[string]bool{"pending": true, "cancelled": true, "completed": true}
	if !validStatuses[status] {
		return models.OrderResponse{}, errors.New("invalid status. Must be pending, completed, or cancelled")
	}

	if err := s.repo.UpdateOrderStatus(id, status); err != nil {
		return models.OrderResponse{}, err
	}

	user, err := s.userClient.ValidateUser(order.UserID)
	orderResponse := models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Products:    order.Products,
		TotalAmount: order.TotalAmount,
		Status:      status,
		CreatedAt:   order.CreatedAt,
		UpdateAt:    time.Now(),
	}
	if err == nil {
		orderResponse.UserName = user.Name
		orderResponse.UserEmail = user.Email
	}
	return orderResponse, nil
}

func (s *OrderService) DeleteOrder(id string) error {
	return s.repo.DeleteOrder(id)
}

func (s *OrderService) GetOrderByID(id string) (models.OrderResponse, error) {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return models.OrderResponse{}, err
	}

	user, err := s.userClient.ValidateUser(order.UserID)
	if err != nil {
		return models.OrderResponse{
			ID:          order.ID,
			Products:    order.Products,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdateAt:    order.UpdateAt,
		}, nil
	}

	return models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		UserName:    user.Name,
		UserEmail:   user.Email,
		Products:    order.Products,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdateAt:    order.UpdateAt,
	}, nil
}
