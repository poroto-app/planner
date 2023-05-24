package entity

import "poroto.app/poroto/planner/internal/domain/models"

type GeoLocationEntity struct {
	Latitude  float64 `firestore:"latitude"`
	Longitude float64 `firestore:"longitude"`
}

func ToGeoLocationEntity(geoLocation models.GeoLocation) GeoLocationEntity {
	return GeoLocationEntity{
		Latitude:  geoLocation.Latitude,
		Longitude: geoLocation.Longitude,
	}
}

func FromGeoLocationEntity(entity GeoLocationEntity) models.GeoLocation {
	return models.GeoLocation{
		Latitude:  entity.Latitude,
		Longitude: entity.Longitude,
	}
}
