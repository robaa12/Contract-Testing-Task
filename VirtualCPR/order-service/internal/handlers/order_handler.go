package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/internal/models"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/internal/service"
)

type OrderServiceInterface interface {
	CreateOrder(req models.CreateOrderRequest) (models.OrderResponse, error)
	GetOrder(id string) (models.OrderResponse, error)
	GetOrderByUserID(userID string) ([]models.OrderResponse, error)
	ListOrders() ([]models.OrderResponse, error)
	UpdateOrderStatus(id string, status string) (models.OrderResponse, error)
	DeleteOrder(id string) error
}

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	orders := router.Group("/api/orders")
	{
		orders.POST("", h.CreateOrder)
		orders.GET("", h.ListOrders)
		orders.GET("/:id", h.GetOrder)
		orders.GET("/user/:userId", h.GetOrderByUser)
		orders.PUT("/:id/status", h.UpdateOrderStatus)
		orders.DELETE("/:id", h.DeleteOrder)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request models.CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrderByUser(c *gin.Context) {
	userId := c.Param("userId")

	orders, err := h.orderService.GetOrderByUserID(userId)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := h.orderService.ListOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var request models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.UpdateOrderStatus(id, request.Status)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	err := h.orderService.DeleteOrder(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}
