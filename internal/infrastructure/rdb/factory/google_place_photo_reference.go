package factory

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewGooglePlacePhotoReferenceFromEntity(
	googlePlacePhotoReferenceEntity entities.GooglePlacePhotoReference,
	googlePlacePhotoAttributionEntities entities.GooglePlacePhotoAttributionSlice,
) models.GooglePlacePhotoReference {
	// HTMLAttributionsを取得
	googlePlacePhotoStrAttributions := array.MapAndFilter(googlePlacePhotoAttributionEntities, func(googlePlacePhotoAttributionEntity *entities.GooglePlacePhotoAttribution) (string, bool) {
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

func NewGooglePlacePhotoReferenceEntityFromGooglePhotoReference(googlePhotoReference models.GooglePlacePhotoReference) entities.GooglePlacePhotoReference {
	return entities.GooglePlacePhotoReference{
		PhotoReference: googlePhotoReference.PhotoReference,
		Width:          googlePhotoReference.Width,
		Height:         googlePhotoReference.Height,
	}
}

func NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePlacePhotoReferences []models.GooglePlacePhotoReference) entities.GooglePlacePhotoReferenceSlice {
	var googlePlacePhotoReferenceEntities entities.GooglePlacePhotoReferenceSlice
	for _, googlePlacePhotoReference := range googlePlacePhotoReferences {
		pr := NewGooglePlacePhotoReferenceEntityFromGooglePhotoReference(googlePlacePhotoReference)
		googlePlacePhotoReferenceEntities = append(googlePlacePhotoReferenceEntities, &pr)
	}
	return googlePlacePhotoReferenceEntities
}
