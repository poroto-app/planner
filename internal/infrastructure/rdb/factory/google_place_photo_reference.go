package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewGooglePlacePhotoReferenceFromEntity(
	googlePlacePhotoReferenceEntity generated.GooglePlacePhotoReference,
	googlePlacePhotoAttributionEntities generated.GooglePlacePhotoAttributionSlice,
) models.GooglePlacePhotoReference {
	// HTMLAttributionsを取得
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

func NewGooglePlacePhotoReferenceEntityFromGooglePhotoReference(googlePhotoReference models.GooglePlacePhotoReference) generated.GooglePlacePhotoReference {
	return generated.GooglePlacePhotoReference{
		PhotoReference: googlePhotoReference.PhotoReference,
		Width:          googlePhotoReference.Width,
		Height:         googlePhotoReference.Height,
	}
}

func NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePlacePhotoReferences []models.GooglePlacePhotoReference) generated.GooglePlacePhotoReferenceSlice {
	var googlePlacePhotoReferenceEntities generated.GooglePlacePhotoReferenceSlice
	for _, googlePlacePhotoReference := range googlePlacePhotoReferences {
		pr := NewGooglePlacePhotoReferenceEntityFromGooglePhotoReference(googlePlacePhotoReference)
		googlePlacePhotoReferenceEntities = append(googlePlacePhotoReferenceEntities, &pr)
	}
	return googlePlacePhotoReferenceEntities
}
