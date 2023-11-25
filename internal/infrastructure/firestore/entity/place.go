package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceEntity struct {
	Id        string  `firestore:"id"`
	Name      string  `firestore:"name"`
	Latitude  float64 `firestore:"latitude"`
	Longitude float64 `firestore:"longitude"`
}

func ToPlaceEntity(place models.Place) PlaceEntity {
	return PlaceEntity{
		Id:        place.Id,
		Name:      place.Name,
		Latitude:  place.Location.Latitude,
		Longitude: place.Location.Longitude,
	}
}

func FromPlaceEntity(entity PlaceEntity) models.Place {
	return models.Place{
		Id:   entity.Id,
		Name: entity.Name,
		Location: models.GeoLocation{
			Latitude:  entity.Latitude,
			Longitude: entity.Longitude,
		},
	}
}
