package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

type GooglePlaceEntity struct {
	PlaceID          string                         `firestore:"place_id"`
	Name             string                         `firestore:"name"`
	Types            []string                       `firestore:"types"`
	Location         GeoLocationEntity              `firestore:"location"`
	PhotoReferences  []string                       `firestore:"photo_references"`
	OpenNow          bool                           `firestore:"open_now"`
	Rating           float32                        `firestore:"rating"`
	UserRatingsTotal int                            `firestore:"user_ratings_total"`
	PriceLevel       int                            `firestore:"price_level"`
	OpeningHours     *GooglePlaceOpeningHoursEntity `firestore:"opening_hours"`
}

func GooglePlaceEntityFromGooglePlace(place models.GooglePlace) GooglePlaceEntity {
	var openingHours *GooglePlaceOpeningHoursEntity
	if place.PlaceDetail != nil {
		o := GooglePlaceOpeningsEntityFromGooglePlaceOpeningPeriod(*place.PlaceDetail.OpeningHours)
		openingHours = &o
	}

	return GooglePlaceEntity{
		PlaceID:          place.PlaceId,
		Name:             place.Name,
		Types:            place.Types,
		PhotoReferences:  place.PhotoReferences,
		OpenNow:          place.OpenNow,
		Rating:           place.Rating,
		UserRatingsTotal: place.UserRatingsTotal,
		PriceLevel:       place.PriceLevel,
		OpeningHours:     openingHours,
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

	var placeDetail *models.GooglePlaceDetail
	if g.OpeningHours != nil {
		placeDetail = &models.GooglePlaceDetail{}
		o := g.OpeningHours.ToGooglePlaceOpeningPeriods()
		placeDetail.OpeningHours = &o
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
