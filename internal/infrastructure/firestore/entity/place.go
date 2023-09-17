package entity

import "poroto.app/poroto/planner/internal/domain/models"

type PlaceEntity struct {
	Id            string            `firestore:"id"`
	GooglePlaceId *string           `firestore:"google_place_id"`
	Name          string            `firestore:"name"`
	Location      GeoLocationEntity `firestore:"location"`
	Thumbnails    []string          `firestore:"thumbnails,omitempty"`
	Photos        []string          `firestore:"photos"`
	// MEMO: Firestoreではuintをサポートしていないため，intにしている
	EstimatedStayDuration int `firestore:"estimated_stay_duration"`
}

func ToPlaceEntity(place models.Place) PlaceEntity {
	return PlaceEntity{
		Id:                    place.Id,
		GooglePlaceId:         place.GooglePlaceId,
		Name:                  place.Name,
		Location:              ToGeoLocationEntity(place.Location),
		Thumbnails:            place.Thumbnails,
		Photos:                place.Photos,
		EstimatedStayDuration: int(place.EstimatedStayDuration),
	}
}

func FromPlaceEntity(entity PlaceEntity) models.Place {
	return models.Place{
		Id:                    entity.Id,
		GooglePlaceId:         entity.GooglePlaceId,
		Name:                  entity.Name,
		Location:              FromGeoLocationEntity(entity.Location),
		Thumbnails:            entity.Thumbnails,
		Photos:                entity.Photos,
		EstimatedStayDuration: uint(entity.EstimatedStayDuration),
	}
}
