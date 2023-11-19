package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlacePhotoReferencesFromPlaceDetail(placeDetail places.PlaceDetail) []models.GooglePlacePhotoReferences {
	var photoReferences []models.GooglePlacePhotoReferences
	for _, photoReference := range placeDetail.Photos {
		photoReferences = append(photoReferences, models.GooglePlacePhotoReferences{
			PhotoReference:   photoReference.PhotoReference,
			Width:            photoReference.Width,
			Height:           photoReference.Height,
			HTMLAttributions: photoReference.HTMLAttributions,
		})
	}

	return photoReferences
}
