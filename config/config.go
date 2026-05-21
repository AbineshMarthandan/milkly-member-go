package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/labstack/gommon/log"
)

type Config struct {
	MongoDB    *mongo.Database
	Port       string
	MongoURL   string
	DBName     string
}

func NewConfig() (*Config, error) {
	// Get configuration from environment variables
	port := getEnv("PORT", "8080")
	mongoURL := getEnv("MONGO_URL", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "milkly-member")

	// Initialize MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Info("Successfully connected to MongoDB")

	database := client.Database(dbName)

	return &Config{
		MongoDB:  database,
		Port:     port,
		MongoURL: mongoURL,
		DBName:   dbName,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}