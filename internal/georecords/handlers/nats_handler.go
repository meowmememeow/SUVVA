package records_handlers

import (
	"context"
	"encoding/json"
	"suvva-geo-ride-service/internal/georecords/dto"
	"suvva-geo-ride-service/internal/georecords/models"
	"suvva-geo-ride-service/internal/logger"

	"github.com/nats-io/nats.go/jetstream"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleRecordNats(msg jetstream.Msg, client *mongo.Client) {

	defer msg.Ack()

	geoRecordsCollection := client.Database("georecords_db").Collection("geo_records")
	var dto dto.InputGeoRecordDto

	if err := json.Unmarshal(msg.Data(), &dto); err != nil {
		logger.ErrorLogger.Printf("Error unmarshalling record: %v", err)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		logger.ErrorLogger.Printf("Invalid ID format")
		return
	}

	record := models.Record{
		ID:       objectID,
		Position: dto.Position,
	}

	_, err = geoRecordsCollection.InsertOne(context.TODO(), record)
	if err != nil {
		logger.ErrorLogger.Printf("Error while inserting geoRecord into database %v", err)
		return
	}

}
