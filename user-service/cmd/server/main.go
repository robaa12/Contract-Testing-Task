package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/user-service/internal/handlers"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/user-service/internal/repository"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/user-service/internal/service"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/user-service/pkg/database"
)

func main() {
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

	userRepo := repository.NewPostgresRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Setup Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register routes
	userHandler.RegisterRoutes(router)

	// Start the server
	port := GetEnvOrDefault("PORT", "8080")
	log.Printf("Server starting on port %s", port)
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
