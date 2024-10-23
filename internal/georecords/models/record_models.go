package models

import (
	shared_models "suvva-geo-ride-service/internal/shared/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Record struct {
	ID       primitive.ObjectID   `bson:"_id"`
	Position shared_models.LngLat `bson:"position"`
	Geozones []string             `bson:"geozones, omitempty"`
}

type InputRecordID struct {
	ID string `json:"id"`
}
