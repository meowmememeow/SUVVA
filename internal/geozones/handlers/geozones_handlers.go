package geozones_handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"net/http"
	"suvva-geo-ride-service/internal/geozones/converters"
	"suvva-geo-ride-service/internal/geozones/dto"
	service "suvva-geo-ride-service/internal/geozones/services"
	"suvva-geo-ride-service/internal/logger"
	shared_models "suvva-geo-ride-service/internal/shared/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateGeozone(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		geozonesCollection := client.Database("geozones_db").Collection("geozones")

		var dto dto.CreateGeozoneDto

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &dto); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(dto.Polyline) == 0 {
			http.Error(w, "Coordinates must not be empty", http.StatusBadRequest)
			return
		}

		geozone := service.CreateGeozoneService(dto)

		_, err = geozonesCollection.InsertOne(context.TODO(), geozone)
		if err != nil {
			http.Error(w, "Error while inserting geozone into database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Geozone created successfully"}`))
	}
}

func GetGeozoneById(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		geozoneID := mux.Vars(r)["geozoneId"]

		objID, err := primitive.ObjectIDFromHex(geozoneID)
		if err != nil {
			http.Error(w, "Invalid geozone ID", http.StatusBadRequest)
			return
		}

		collection := client.Database("geozones_db").Collection("geozones")
		filter := bson.M{"_id": objID, "deleted": bson.M{"$ne": true}}

		var geozone shared_models.Geozone
		err = collection.FindOne(context.Background(), filter).Decode(&geozone)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, "Geozone not found", http.StatusNotFound)
			} else {
				logger.ErrorLogger.Printf("Error finding geozone: %v", err)
				http.Error(w, "Error finding geozone", http.StatusInternalServerError)
			}
			return
		}

		geozoneDto := converters.GeozoneDtoFronGeozone(geozone)

		data, _ := json.Marshal(geozoneDto)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func UpdateGeozone(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		geozoneID := mux.Vars(r)["geozoneId"]

		objID, err := primitive.ObjectIDFromHex(geozoneID)
		if err != nil {
			http.Error(w, "Invalid geozone ID", http.StatusBadRequest)
			return
		}

		var dto dto.UpdateGeozoneDto

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &dto); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := json.Marshal(service.UpdateGeozoneService(objID, dto, client))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Geozone updated successfully"}`))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func DeleteGeozone(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		geozoneID := mux.Vars(r)["geozoneId"]

		objID, err := primitive.ObjectIDFromHex(geozoneID)
		if err != nil {
			http.Error(w, "Invalid geozone ID", http.StatusBadRequest)
			return
		}

		now := time.Now()
		update := bson.M{
			"$set": bson.M{
				"deleted":   true,
				"deletedAt": now,
				"updatedAt": now,
			},
		}

		collection := client.Database("geozones_db").Collection("geozones")
		result := collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objID, "deleted": bson.M{"$ne": true}}, update)
		if result.Err() != nil {
			logger.ErrorLogger.Printf("Error deleting geozone: %v", result.Err())
			http.Error(w, "Error deleting geozone", http.StatusInternalServerError)
			return
		}

		var deletedGeozone shared_models.Geozone
		if err := result.Decode(&deletedGeozone); err != nil {
			logger.ErrorLogger.Printf("Error decoding deleted geozone: %v", err)
			http.Error(w, "Error decoding deleted geozone", http.StatusInternalServerError)
			return
		}

		data, _ := json.Marshal(deletedGeozone)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
