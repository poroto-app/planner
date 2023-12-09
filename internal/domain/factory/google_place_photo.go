package factory

import (
	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlacePhotoReferencesFromPlaceDetail(placeDetail places.PlaceDetail) []models.GooglePlacePhotoReference {
	var photoReferences []models.GooglePlacePhotoReference
	for _, photo := range placeDetail.Photos {
		photoReferences = append(photoReferences, GooglePlacePhotoReferenceFromPhoto(photo))
	}

	return photoReferences
}

func GooglePlacePhotoReferenceFromPhoto(photo maps.Photo) models.GooglePlacePhotoReference {
	return models.GooglePlacePhotoReference{
		PhotoReference:   photo.PhotoReference,
		Width:            photo.Width,
		Height:           photo.Height,
		HTMLAttributions: photo.HTMLAttributions,
	}
}
