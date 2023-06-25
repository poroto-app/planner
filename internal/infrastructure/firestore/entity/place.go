package entity

import "poroto.app/poroto/planner/internal/domain/models"

type PlaceEntity struct {
	Id        string            `json:"id"`
	Name      string            `firestore:"name"`
	Location  GeoLocationEntity `firestore:"location"`
	Thumbnail *string           `firestore:"thumbnail"`
	Photos    []string          `firestore:"photos"`
	// MEMO: Firestoreではuintをサポートしていないため，intにしている
	EstimatedStayDuration int `firestore:"estimated_stay_duration"`
}

func ToPlaceEntity(place models.Place) PlaceEntity {
	return PlaceEntity{
		Id:                    place.Id,
		Name:                  place.Name,
		Location:              ToGeoLocationEntity(place.Location),
		Thumbnail:             place.Thumbnail,
		Photos:                place.Photos,
		EstimatedStayDuration: int(place.EstimatedStayDuration),
	}
}

func FromPlaceEntity(entity PlaceEntity) models.Place {
	return models.Place{
		Id:                    entity.Id,
		Name:                  entity.Name,
		Location:              FromGeoLocationEntity(entity.Location),
		Thumbnail:             entity.Thumbnail,
		Photos:                entity.Photos,
		EstimatedStayDuration: uint(entity.EstimatedStayDuration),
	}
}
