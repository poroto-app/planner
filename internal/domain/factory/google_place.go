package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlaceFromPlaceEntity(place googleplaces.Place, photos *[]models.GooglePlacePhoto) models.GooglePlace {
	var placeDetail *models.GooglePlaceDetail
	if place.PlaceDetail != nil {
		d := GooglePlaceDetailFromPlaceDetailEntity(*place.PlaceDetail)
		placeDetail = &d
	}

	return models.GooglePlace{
		PlaceId: place.PlaceID,
		Name:    place.Name,
		Types:   place.Types,
		Location: models.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
		PhotoReferences:  place.PhotoReferences,
		OpenNow:          place.OpenNow,
		Rating:           place.Rating,
		UserRatingsTotal: place.UserRatingsTotal,
		PriceLevel:       place.PriceLevel,
		Photos:           photos,
		PlaceDetail:      placeDetail,
	}
}

func PlaceEntityFromGooglePlace(place models.GooglePlace) googleplaces.Place {
	return googleplaces.Place{
		PlaceID: place.PlaceId,
		Name:    place.Name,
		Types:   place.Types,
		Location: googleplaces.Location{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
		PhotoReferences:  place.PhotoReferences,
		OpenNow:          place.OpenNow,
		Rating:           place.Rating,
		UserRatingsTotal: place.UserRatingsTotal,
		PriceLevel:       place.PriceLevel,
	}
}
