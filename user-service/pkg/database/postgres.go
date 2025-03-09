package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	maxRetries := 30
	retryInterval := 5 * time.Second
	// Open a connection to the database
	var db *sql.DB
	var err error
	for attempt := 1; attempt < maxRetries; attempt++ {
		log.Printf("Trying to connect to the database. Attempt %d", attempt)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Println("Failed to open database connection")
			time.Sleep(retryInterval)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Printf("Connected to database!")
			break
		}
		log.Printf("Failed to ping database: %v. Retrying in %v...", err, retryInterval)
		db.Close()
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

func GetConfigFromEnv() Config {
	return Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvAsIntOrDDefault("DB_PORT", 5432),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "password"),
		DBName:   getEnvOrDefault("DB_NAME", "user_service"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsIntOrDDefault(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
