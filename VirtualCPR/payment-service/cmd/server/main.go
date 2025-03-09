package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/internal/handlers"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/internal/repository"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/internal/service"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/payment-service/pkg/database"
	"github.com/stripe/stripe-go/v81"
)

func main() {
	// Initialize Stripe
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		stripeKey = "pk_test_51QzteqEN3C714OAmopACj4peCAlnLnU5o4LSQlaMg0m3q5XV0GwZ1vVbHTh2YBktcIVFN2us9vevw8lsPuCPz1dk00Eu1o6Rb7" // Default test key for development
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	stripe.SetHTTPClient(httpClient)
	stripe.Key = stripeKey

	dbConfig := database.GetConfigFromEnv()

	// Connect to database
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup database schema
	if err := database.SetupSchema(db); err != nil {
		log.Fatalf("Failed to setup database schema: %v", err)
	}

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Setup Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register routes
	paymentHandler.RegisterRoutes(router)

	// Start the server
	port := GetEnvOrDefault("PORT", "8080")
	log.Printf("Payment service starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
