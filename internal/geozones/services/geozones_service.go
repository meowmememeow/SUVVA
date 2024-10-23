package service

import (
	"context"
	"suvva-geo-ride-service/internal/geozones/dto"
	"suvva-geo-ride-service/internal/helpers"
	shared_models "suvva-geo-ride-service/internal/shared/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateGeozoneService(dto dto.CreateGeozoneDto) (geozone shared_models.Geozone) {
	var lineString shared_models.LineString
	for _, p := range dto.Polyline {
		point := helpers.CoordinateFromLngLat(&p)
		lineString = append(lineString, point)
	}
	coordinates := shared_models.Polygon{lineString}

	geoPolygon := shared_models.GeoJSONPolygon{
		Type:        "Polygon",
		Coordinates: coordinates,
	}

	geozone = shared_models.Geozone{
		ID:        primitive.NewObjectID(),
		Name:      dto.Name,
		Polygon:   geoPolygon,
		Deleted:   false,
		UpdatedAt: time.Now(),
	}
	return geozone
}

func UpdateGeozoneService(objID primitive.ObjectID, dto dto.UpdateGeozoneDto, client *mongo.Client) (updatedGeozone shared_models.Geozone) {

	geozonesCollection := client.Database("geozones_db").Collection("geozones")

	update := bson.M{"$set": bson.M{"updatedAt": time.Now()}}

	updatedFields := bson.M{
		"updatedAt": time.Now(),
	}
	if dto.Name != nil {
		updatedFields["updated name"] = *dto.Name
	}

	if dto.Polyline != nil {
		var lineString shared_models.LineString
		for _, p := range dto.Polyline {
			point := helpers.CoordinateFromLngLat(&p)
			lineString = append(lineString, point)
		}
		coordinates := shared_models.Polygon{lineString}

		geoPolygon := shared_models.GeoJSONPolygon{
			Type:        "Polygon",
			Coordinates: coordinates,
		}
		updatedFields["updated polygon"] = geoPolygon
	}

	update = bson.M{"$set": updatedFields}

	result := geozonesCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": objID, "deleted": bson.M{"$ne": true}}, update)
	if result.Err() != nil {
		return
	}

	if err := result.Decode(&updatedGeozone); err != nil {
		return
	}
	return updatedGeozone
}
