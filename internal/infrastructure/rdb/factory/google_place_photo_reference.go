package factory

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewGooglePlacePhotoReferenceFromEntity(
	googlePlacePhotoReferenceEntity generated.GooglePlacePhotoReference,
	googlePlacePhotoAttributionEntities generated.GooglePlacePhotoAttributionSlice,
) models.GooglePlacePhotoReference {
	// HTMLAttributionsを取得
	googlePlacePhotoStrAttributions := array.MapAndFilter(googlePlacePhotoAttributionEntities, func(googlePlacePhotoAttributionEntity *generated.GooglePlacePhotoAttribution) (string, bool) {
		if googlePlacePhotoAttributionEntity == nil {
			return "", false
		}

		// PhotoReferenceが一致するものだけを抽出
		if googlePlacePhotoAttributionEntity.PhotoReference != googlePlacePhotoReferenceEntity.PhotoReference {
			return "", false
		}

		return googlePlacePhotoAttributionEntity.HTMLAttribution, true
	})

	return models.GooglePlacePhotoReference{
		PhotoReference:   googlePlacePhotoReferenceEntity.PhotoReference,
		Width:            googlePlacePhotoReferenceEntity.Width,
		Height:           googlePlacePhotoReferenceEntity.Height,
		HTMLAttributions: googlePlacePhotoStrAttributions,
	}
}

func NewGooglePlacePhotoReferenceEntityFromGooglePhotoReference(googlePhotoReference models.GooglePlacePhotoReference, googlePlaceId string) generated.GooglePlacePhotoReference {
	return generated.GooglePlacePhotoReference{
		PhotoReference: googlePhotoReference.PhotoReference,
		GooglePlaceID:  googlePlaceId,
		Width:          googlePhotoReference.Width,
		Height:         googlePhotoReference.Height,
	}
}

func NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePlacePhotoReferences []models.GooglePlacePhotoReference, googlePlaceId string) generated.GooglePlacePhotoReferenceSlice {
	var googlePlacePhotoReferenceEntities generated.GooglePlacePhotoReferenceSlice
	for _, googlePlacePhotoReference := range googlePlacePhotoReferences {
		pr := NewGooglePlacePhotoReferenceEntityFromGooglePhotoReference(googlePlacePhotoReference, googlePlaceId)
		googlePlacePhotoReferenceEntities = append(googlePlacePhotoReferenceEntities, &pr)
	}
	return googlePlacePhotoReferenceEntities
}
