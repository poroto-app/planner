package factory

import (
	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlacePhotoReferencesFromPlaceDetail(placeDetail places.PlaceDetail) []models.GooglePlacePhotoReference {
	photoReferences := make([]models.GooglePlacePhotoReference, len(placeDetail.Photos))
	for i, photo := range placeDetail.Photos {
		photoReferences[i] = GooglePlacePhotoReferenceFromPhoto(photo)
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
