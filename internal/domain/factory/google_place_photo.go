package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlacePhotosFromPlaceDetail(placeDetail places.PlaceDetail) []models.GooglePlacePhoto {
	var photos []models.GooglePlacePhoto
	for _, photo := range placeDetail.Photos {
		photos = append(photos, models.GooglePlacePhoto{
			PhotoReference:   photo.PhotoReference,
			Width:            photo.Width,
			Height:           photo.Height,
			HTMLAttributions: photo.HTMLAttributions,
		})
	}

	return photos
}
