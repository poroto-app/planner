package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewGooglePlacePhotoReferenceFromEntity(
	googlePlacePhotoReferenceEntity entities.GooglePlacePhotoReference,
	googlePlacePhotoAttributionEntities entities.GooglePlacePhotoAttributionSlice,
) models.GooglePlacePhotoReference {
	var googlePlacePhotoStrAttributions []string
	for _, googlePlacePhotoAttributionEntity := range googlePlacePhotoAttributionEntities {
		if googlePlacePhotoAttributionEntity == nil {
			continue
		}

		// PhotoReferenceが一致するものだけを抽出
		if googlePlacePhotoAttributionEntity.PhotoReference != googlePlacePhotoReferenceEntity.PhotoReference {
			continue
		}

		googlePlacePhotoStrAttributions = append(googlePlacePhotoStrAttributions, googlePlacePhotoAttributionEntity.HTMLAttribution)
	}

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
