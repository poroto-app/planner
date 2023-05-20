package entity

import "poroto.app/poroto/planner/internal/domain/models"

type PlaceEntity struct {
	Name                  string            `firestore:"name"`
	Location              GeoLocationEntity `firestore:"location"`
	Thumbnail             *string           `firestore:"thumbnail"`
	Photos                []string          `firestore:"photos"`
	EstimatedStayDuration uint              `firestore:"estimated_stay_duration"`
}

func ToPlaceEntity(place models.Place) PlaceEntity {
	return PlaceEntity{
		Name:                  place.Name,
		Location:              ToGeoLocationEntity(place.Location),
		Thumbnail:             place.Thumbnail,
		Photos:                place.Photos,
		EstimatedStayDuration: place.EstimatedStayDuration,
	}
}

func FromPlaceEntity(entity PlaceEntity) models.Place {
	return models.Place{
		Name:                  entity.Name,
		Location:              FromGeoLocationEntity(entity.Location),
		Thumbnail:             entity.Thumbnail,
		Photos:                entity.Photos,
		EstimatedStayDuration: entity.EstimatedStayDuration,
	}
}
