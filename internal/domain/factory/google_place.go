package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

func GooglePlaceFromPlaceEntity(place googleplaces.Place, imageEntities []entity.ImageEntity, reviewEntities []entity.GooglePlaceReviewEntity) models.GooglePlace {
	var images *[]models.Image
	if len(imageEntities) == 0 {
		images = nil
	} else {
		images = new([]models.Image)
		for _, imageEntity := range imageEntities {
			*images = append(*images, entity.FromImageEntity(imageEntity))
		}
	}

	var reviews *[]models.GooglePlaceReview
	if len(reviewEntities) == 0 {
		reviews = nil
	} else {
		reviews = new([]models.GooglePlaceReview)
		for _, reviewEntity := range reviewEntities {
			*reviews = append(*reviews, entity.FromGooglePlaceReviewEntity(reviewEntity))
		}
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
		Images:           images,
		Reviews:          reviews,
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
	}
}
