package config

import (
	"errors"
	"os"
)

type Config struct {
	MongoURI    string
	NatsURL     string
	StreamName  string
	SubjectName string
	DurableName string
}

var ConfigInstance *Config

func LoadConfig() error {

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return errors.New("missing required environment variable: MONGO_URI")
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return errors.New("missing required environment variable: NATS_URL")
	}

	streamName := os.Getenv("STREAM_NAME")
	if streamName == "" {
		streamName = "RECORDS_STREAM"
	}

	subjectName := os.Getenv("SUBJECT_NAME")
	if subjectName == "" {
		subjectName = "records"
	}

	durableName := os.Getenv("DURABLE_NAME")
	if durableName == "" {
		durableName = "CONS"
	}

	ConfigInstance = &Config{
		MongoURI:    mongoURI,
		NatsURL:     natsURL,
		StreamName:  streamName,
		SubjectName: subjectName,
		DurableName: durableName,
	}

	return nil
}
