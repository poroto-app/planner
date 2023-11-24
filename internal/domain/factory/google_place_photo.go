package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlacePhotoReferencesFromPlaceDetail(placeDetail places.PlaceDetail) []models.GooglePlacePhotoReference {
	var photoReferences []models.GooglePlacePhotoReference
	for _, photoReference := range placeDetail.Photos {
		photoReferences = append(photoReferences, models.GooglePlacePhotoReference{
			PhotoReference:   photoReference.PhotoReference,
			Width:            photoReference.Width,
			Height:           photoReference.Height,
			HTMLAttributions: photoReference.HTMLAttributions,
		})
	}

	return photoReferences
}
