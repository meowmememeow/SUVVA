package helpers

import (
	"suvva-geo-ride-service/internal/geozones/dto"
	shared_models "suvva-geo-ride-service/internal/shared/models"
)

func CoordinateFromLngLat(lngLat *dto.LngLatDto) [2]float64 {
	return [2]float64{lngLat.Longitude, lngLat.Latitude}
}

func LngLatFromCoordinate(coord *shared_models.Point) dto.LngLatDto {
	dto := dto.LngLatDto{
		Longitude: coord[0],
		Latitude:  coord[1],
	}
	return dto
}
