package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewGooglePlacePhotoAttributionSliceFromPhotoReference(googlePlacePhotoReference models.GooglePlacePhotoReference, googlePlaceId string) entities.GooglePlacePhotoAttributionSlice {
	photoAttributions := make(entities.GooglePlacePhotoAttributionSlice, len(googlePlacePhotoReference.HTMLAttributions))
	for i, attribution := range googlePlacePhotoReference.HTMLAttributions {
		photoAttributions[i] = &entities.GooglePlacePhotoAttribution{
			ID:              uuid.New().String(),
			GooglePlaceID:   googlePlaceId,
			PhotoReference:  googlePlacePhotoReference.PhotoReference,
			HTMLAttribution: attribution,
		}
	}
	return photoAttributions
}
