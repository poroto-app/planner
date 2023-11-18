package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

type GooglePlaceEntity struct {
	PlaceID          string            `firestore:"place_id"`
	Name             string            `firestore:"name"`
	Types            []string          `firestore:"types"`
	Location         GeoLocationEntity `firestore:"location"`
	PhotoReferences  []string          `firestore:"photo_references"`
	OpenNow          bool              `firestore:"open_now"`
	Rating           float32           `firestore:"rating"`
	UserRatingsTotal int               `firestore:"user_ratings_total"`
	PriceLevel       int               `firestore:"price_level"`
}

func GooglePlaceEntityFromGooglePlace(place models.GooglePlace) GooglePlaceEntity {
	return GooglePlaceEntity{
		PlaceID:          place.PlaceId,
		Name:             place.Name,
		Types:            place.Types,
		PhotoReferences:  place.PhotoReferences,
		OpenNow:          place.OpenNow,
		Rating:           place.Rating,
		UserRatingsTotal: place.UserRatingsTotal,
		PriceLevel:       place.PriceLevel,
		Location: GeoLocationEntity{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
	}
}

func (g GooglePlaceEntity) ToGooglePlace(images *[]models.Image) models.GooglePlace {
	location := models.GeoLocation{
		Latitude:  g.Location.Latitude,
		Longitude: g.Location.Longitude,
	}

	// TODO: Place Detailを復元する
	return models.GooglePlace{
		PlaceId:          g.PlaceID,
		Name:             g.Name,
		Types:            g.Types,
		Location:         location,
		PhotoReferences:  g.PhotoReferences,
		OpenNow:          g.OpenNow,
		Rating:           g.Rating,
		UserRatingsTotal: g.UserRatingsTotal,
		PriceLevel:       g.PriceLevel,
		Images:           images,
	}
}
