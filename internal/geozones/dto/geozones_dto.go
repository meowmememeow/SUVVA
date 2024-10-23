package dto

type LngLatDto struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}

type CreateGeozoneDto struct {
	Name     string      `json:"name" validate:"required"`
	Polyline []LngLatDto `json:"polyline" validate:"required"`
}

type UpdateGeozoneDto struct {
	Name     *string     `json:"name"  validate:"required"`
	Polyline []LngLatDto `json:"polyline"  validate:"required"`
}

type GeozoneDto struct {
	ID       string      `json:"id" validate:"required"`
	Name     string      `json:"name" validate:"required"`
	Polyline []LngLatDto `json:"polyline" validate:"required"`
}
