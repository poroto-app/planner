package factory

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewGooglePlaceDetailFromGooglePlaceEntity(googlePlaceEntity entities.GooglePlace) (*models.GooglePlaceDetail, error) {
	if googlePlaceEntity.R == nil {
		return nil, fmt.Errorf("googlePlaceEntity.R is nil")
	}

	googlePlaceReviews := array.MapAndFilter(googlePlaceEntity.R.GetGooglePlaceReviews(), func(googlePlaceReviewEntity *entities.GooglePlaceReview) (models.GooglePlaceReview, bool) {
		if googlePlaceReviewEntity == nil {
			return models.GooglePlaceReview{}, false
		}

		if googlePlaceReviewEntity.GooglePlaceID != googlePlaceEntity.GooglePlaceID {
			return models.GooglePlaceReview{}, false
		}

		return NewGooglePlaceReviewFromEntity(*googlePlaceReviewEntity), true
	})

	googlePlaceOpeningPeriods := array.MapAndFilter(googlePlaceEntity.R.GetGooglePlaceOpeningPeriods(), func(googlePlaceOpeningPeriodEntity *entities.GooglePlaceOpeningPeriod) (models.GooglePlaceOpeningPeriod, bool) {
		if googlePlaceOpeningPeriodEntity == nil {
			return models.GooglePlaceOpeningPeriod{}, false
		}

		if googlePlaceOpeningPeriodEntity.GooglePlaceID != googlePlaceEntity.GooglePlaceID {
			return models.GooglePlaceOpeningPeriod{}, false
		}

		return NewGooglePlaceOpeningPeriodFromEntity(*googlePlaceOpeningPeriodEntity), true
	})

	googlePlacePhotoReferenceEntities := array.Filter(googlePlaceEntity.R.GetGooglePlacePhotoReferences(), func(googlePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) bool {
		if googlePlacePhotoReferenceEntity == nil {
			return false
		}

		if googlePlacePhotoReferenceEntity.GooglePlaceID != googlePlaceEntity.GooglePlaceID {
			return false
		}

		return true
	})
	googlePlacePhotoReferences := array.MapAndFilter(googlePlacePhotoReferenceEntities, func(googlePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) (models.GooglePlacePhotoReference, bool) {
		if googlePlacePhotoReferenceEntity == nil {
			return models.GooglePlacePhotoReference{}, false
		}

		return NewGooglePlacePhotoReferenceFromEntity(*googlePlacePhotoReferenceEntity, googlePlaceEntity.R.GetGooglePlacePhotoAttributions()), true
	})

	if len(googlePlaceReviews) == 0 && len(googlePlaceOpeningPeriods) == 0 && len(googlePlacePhotoReferenceEntities) == 0 {
		return nil, nil
	}

	return &models.GooglePlaceDetail{
		Reviews:         googlePlaceReviews,
		PhotoReferences: googlePlacePhotoReferences,
		OpeningHours: &models.GooglePlaceOpeningHours{
			Periods: googlePlaceOpeningPeriods,
		},
	}, nil
}
