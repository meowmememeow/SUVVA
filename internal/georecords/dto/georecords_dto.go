package dto // в отдельную папку дто положить

import (
	"time"

	shared_models "suvva-geo-ride-service/internal/shared/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OutputGeoRecordDto struct {
	ID        primitive.ObjectID   `bson:"_id" validate:"required"`
	Position  shared_models.LngLat `bson:"position" validate:"required"`
	Geozones  []string             `bson:"geozones" validate:"required"`
	CaptureAt time.Time            `bson:"captureAt" validate:"required"`
}

type InputGeoRecordDto struct {
	ID        string               `json:"id" validate:"required"`
	Position  shared_models.LngLat `json:"position" validate:"required"`
	CaptureAt string               `json:"captureAt" validate:"required"`
}
