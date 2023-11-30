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

func NewPlaceEntityFromPlace(place models.Place) PlaceEntity {
	return PlaceEntity{
		Id:        place.Id,
		Name:      place.Name,
		Latitude:  place.Location.Latitude,
		Longitude: place.Location.Longitude,
	}
}

func (p PlaceEntity) ToPlace() models.Place {
	return models.Place{
		Id:   p.Id,
		Name: p.Name,
		Location: models.GeoLocation{
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
		},
	}
}
