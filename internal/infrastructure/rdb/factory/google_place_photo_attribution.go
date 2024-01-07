package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewGooglePlacePhotoAttributionSliceFromPhotoReference(googlePlacePhotoReference models.GooglePlacePhotoReference, googlePlaceId string) generated.GooglePlacePhotoAttributionSlice {
	photoAttributions := make(generated.GooglePlacePhotoAttributionSlice, len(googlePlacePhotoReference.HTMLAttributions))
	for i, attribution := range googlePlacePhotoReference.HTMLAttributions {
		photoAttributions[i] = &generated.GooglePlacePhotoAttribution{
			ID:              uuid.New().String(),
			GooglePlaceID:   googlePlaceId,
			PhotoReference:  googlePlacePhotoReference.PhotoReference,
			HTMLAttribution: attribution,
		}
	}
	return photoAttributions
}
