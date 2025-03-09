package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/internal/handlers"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/internal/repository"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/internal/service"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/pkg/client"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/order-service/pkg/database"
)

func main() {
	dbConfig := database.GetConfigFromEnv()

	// Connect to the database
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.SetupSchema(db); err != nil {
		log.Fatalf("Failed to setup schema: %v", err)
	}

	userServiceURL := getEnvOrDefault("USER_SERVICE_URL", "http://localhost:8080")
	userTimeout := 5
	userClient := client.NewHttpUserClient(userServiceURL, userTimeout)

	orderRepo := repository.NewPostgresOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, userClient)
	orderHandler := handlers.NewOrderHandler(*orderService)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	orderHandler.RegisterRoutes(router)

	port := getEnvOrDefault("PORT", "8081")
	log.Printf("Order service starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
