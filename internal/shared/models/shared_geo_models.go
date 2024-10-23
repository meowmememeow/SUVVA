package shared_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LngLat struct {
	Longitude float64 `bson:"lng"`
	Latitude  float64 `bson:"lat"`
}

type Point [2]float64
type LineString []Point
type Polygon []LineString

type GeoJSONPolygon struct {
	Type        string  `bson:"type"`
	Coordinates Polygon `bson:"coordinates"`
}

type Geozone struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Polygon   GeoJSONPolygon     `bson:"polygon"`
	Deleted   bool               `bson:"deleted,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
}
