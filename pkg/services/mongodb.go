package services

import (
	"context"
	"suvva-geo-ride-service/internal/config"
	"suvva-geo-ride-service/internal/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClient() *mongo.Client {

	clientOptions := options.Client().ApplyURI(config.ConfigInstance.MongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	URI := config.ConfigInstance.MongoURI

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		logger.InfoLogger.Println("URI:")
		logger.ErrorLogger.Print(URI)
		logger.ErrorLogger.Fatalf("Failed to ping MongoDB: %v", err)
	}

	logger.InfoLogger.Println("Connected to MongoDB")
	return client
}
