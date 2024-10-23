package services

import (
	"context"
	shared_models "suvva-geo-ride-service/internal/shared/models"

	"suvva-geo-ride-service/internal/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindMatchingGeozones(client *mongo.Client, point shared_models.LngLat) ([]primitive.ObjectID, error) {

	GeozoneCollection := client.Database("geozones_db").Collection("geozones")

	pointGeoJSON := bson.M{
		"type":        "Point",
		"coordinates": []float64{point.Longitude, point.Latitude},
	}

	cursor, err := GeozoneCollection.Find(context.TODO(), bson.M{
		"polygon": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": pointGeoJSON,
			},
		},
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error finding geozones: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var geozones []primitive.ObjectID
	for cursor.Next(context.TODO()) {
		var geozone shared_models.Geozone
		if err := cursor.Decode(&geozone); err != nil {
			logger.ErrorLogger.Printf("Error decoding geozone: %v", err)
			continue
		}
		geozones = append(geozones, geozone.ID)
	}
	if err := cursor.Err(); err != nil {
		logger.ErrorLogger.Printf("Cursor error: %v", err)
		return nil, err
	}

	return geozones, nil
}
