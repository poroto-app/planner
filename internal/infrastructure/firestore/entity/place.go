package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceEntity struct {
	Id            string    `firestore:"id"`
	Name          string    `firestore:"name"`
	GooglePlaceId string    `firestore:"google_place_id"`
	Latitude      float64   `firestore:"latitude"`
	Longitude     float64   `firestore:"longitude"`
	GeoHash       string    `firestore:"geohash"`
	CreatedAt     time.Time `firestore:"created_at,serverTimestamp,omitempty"`
	UpdatedAt     time.Time `firestore:"updated_at,omitempty"`
	LikeCount     uint      `firestore:"like_count"`
}

func NewPlaceEntityFromPlace(place models.Place) PlaceEntity {
	return PlaceEntity{
		Id:            place.Id,
		Name:          place.Name,
		GooglePlaceId: place.Google.PlaceId,
		Latitude:      place.Location.Latitude,
		Longitude:     place.Location.Longitude,
		GeoHash:       place.Location.GeoHash(),
		UpdatedAt:     time.Now(),
	}
}

func NewPlaceEntityFromGooglePlace(placeId string, googlePlace models.GooglePlace) PlaceEntity {
	return PlaceEntity{
		Id:            placeId,
		Name:          googlePlace.Name,
		GooglePlaceId: googlePlace.PlaceId,
		Latitude:      googlePlace.Location.Latitude,
		Longitude:     googlePlace.Location.Longitude,
		GeoHash:       googlePlace.Location.GeoHash(),
		UpdatedAt:     time.Now(),
	}
}

func (p PlaceEntity) ToPlace(googlePlace models.GooglePlace) models.Place {
	return models.Place{
		Id:     p.Id,
		Name:   p.Name,
		Google: googlePlace,
		Location: models.GeoLocation{
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
		},
	}
}
