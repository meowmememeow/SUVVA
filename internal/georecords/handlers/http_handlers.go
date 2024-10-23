package records_handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"suvva-geo-ride-service/internal/georecords/dto"
	"suvva-geo-ride-service/internal/georecords/models"
	"suvva-geo-ride-service/internal/georecords/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateGeoRecord(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		geoRecordsCollection := client.Database("georecords_db").Collection("geo_records")

		var dto dto.InputGeoRecordDto

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

		objectID, err := primitive.ObjectIDFromHex(dto.ID)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		record := models.Record{
			ID:       objectID,
			Position: dto.Position,
		}

		_, err = geoRecordsCollection.InsertOne(context.TODO(), record)
		if err != nil {
			http.Error(w, "Error while inserting geoRecord into database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "GeoRecord created successfully"}`))
	}
}

func UpdateRecordGeozones(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		geoRecordsCollection := client.Database("georecords_db").Collection("geo_records")

		var inputID models.InputRecordID

		if err := json.NewDecoder(r.Body).Decode(&inputID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		objectID, err := primitive.ObjectIDFromHex(inputID.ID)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var record models.Record
		if err := geoRecordsCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&record); err != nil {
			http.Error(w, "Record not found", http.StatusNotFound)
			return
		}

		geozones, err := services.FindMatchingGeozones(client, record.Position)
		if err != nil {
			log.Printf("Error finding geozones: %v", err)
			http.Error(w, "Error while checking geozones", http.StatusInternalServerError)
			return
		}

		_, err = geoRecordsCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{
			"$set": bson.M{"geozones": geozones},
		})
		if err != nil {
			log.Printf("Error updating record with geozones: %v", err)
			http.Error(w, "Error while updating record with geozones", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"geozones": geozones,
		}
		json.NewEncoder(w).Encode(response)
	}
}
