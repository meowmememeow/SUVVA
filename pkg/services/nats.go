package services

import (
	"context"
	"suvva-geo-ride-service/internal/config"
	records_handlers "suvva-geo-ride-service/internal/georecords/handlers"
	"suvva-geo-ride-service/internal/logger"

	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetNatsConnection() (*nats.Conn, error) {
	natsURL := config.ConfigInstance.NatsURL

	//logger.InfoLogger.Printf("Connecting to NATS at %s", natsURL)
	nc, err := nats.Connect(natsURL)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error connecting to NATS: %v", err)
		return nil, err
	}
	logger.InfoLogger.Println("Connected to NATS")
	return nc, nil
}

func GetJetStreamContext(nc *nats.Conn) jetstream.JetStream {
	js, err := jetstream.New(nc)

	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to create JetStream context: %v", err)
	}
	return js
}

func SetupSubscriptions(js jetstream.JetStream, nc *nats.Conn, client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     config.ConfigInstance.StreamName,
		Subjects: []string{config.ConfigInstance.SubjectName},
	})

	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to create stream: %v", err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   config.ConfigInstance.DurableName,
		AckPolicy: jetstream.AckExplicitPolicy,
	})

	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to create consumer: %v", err)
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		records_handlers.HandleRecordNats(msg, client)
	})

	if err != nil {
		logger.ErrorLogger.Fatalf("Error subscribing to NATS topic: %v", err)
	}

	logger.InfoLogger.Println("Subscriptions setup complete")
}
