package converters

import (
	"suvva-geo-ride-service/internal/geozones/dto"
	"suvva-geo-ride-service/internal/helpers"
	shared_models "suvva-geo-ride-service/internal/shared/models"
)

func GeozoneDtoFronGeozone(geozone shared_models.Geozone) (geozoneDto dto.GeozoneDto) {
	var polyline []dto.LngLatDto

	for _, point := range geozone.Polygon.Coordinates[0] {
		polyline = append(polyline, helpers.LngLatFromCoordinate(&point))
	}

	geozoneDto = dto.GeozoneDto{
		ID:       geozone.ID.Hex(),
		Name:     geozone.Name,
		Polyline: polyline,
	}

	return geozoneDto
}
