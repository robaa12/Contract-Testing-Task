package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/internal/models"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/internal/service"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/payments", h.CreatePayment)
	router.GET("/payments/:id", h.GetPaymentByID)
	router.GET("/payments/user/:user_id", h.ListPaymentsByUserID)
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var request models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentService.CreatePayment(request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, payment)
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.paymentService.GetPaymentByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, payment)
}

func (h *PaymentHandler) ListPaymentsByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	payments, err := h.paymentService.ListPaymentsByUserID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, payments)
}
