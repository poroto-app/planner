package factory

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewGooglePlaceDetailFromGooglePlaceEntity(
	googlePlaceReviewSlice entities.GooglePlaceReviewSlice,
	googlePlaceOpeningPeriodSlice entities.GooglePlaceOpeningPeriodSlice,
	googlePlacePhotoReferenceSlice entities.GooglePlacePhotoReferenceSlice,
	googlePlacePhotoAttributionSlice entities.GooglePlacePhotoAttributionSlice,
	googlePlaceId string,
) (*models.GooglePlaceDetail, error) {
	googlePlaceReviews := array.MapAndFilter(googlePlaceReviewSlice, func(googlePlaceReviewEntity *entities.GooglePlaceReview) (models.GooglePlaceReview, bool) {
		if googlePlaceReviewEntity == nil {
			return models.GooglePlaceReview{}, false
		}

		if googlePlaceReviewEntity.GooglePlaceID != googlePlaceId {
			return models.GooglePlaceReview{}, false
		}

		return NewGooglePlaceReviewFromEntity(*googlePlaceReviewEntity), true
	})

	googlePlaceOpeningPeriods := array.MapAndFilter(googlePlaceOpeningPeriodSlice, func(googlePlaceOpeningPeriodEntity *entities.GooglePlaceOpeningPeriod) (models.GooglePlaceOpeningPeriod, bool) {
		if googlePlaceOpeningPeriodEntity == nil {
			return models.GooglePlaceOpeningPeriod{}, false
		}

		if googlePlaceOpeningPeriodEntity.GooglePlaceID != googlePlaceId {
			return models.GooglePlaceOpeningPeriod{}, false
		}

		return NewGooglePlaceOpeningPeriodFromEntity(*googlePlaceOpeningPeriodEntity), true
	})

	googlePlacePhotoReferenceEntities := array.Filter(googlePlacePhotoReferenceSlice, func(googlePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) bool {
		if googlePlacePhotoReferenceEntity == nil {
			return false
		}

		if googlePlacePhotoReferenceEntity.GooglePlaceID != googlePlaceId {
			return false
		}

		return true
	})
	googlePlacePhotoReferences := array.MapAndFilter(googlePlacePhotoReferenceEntities, func(googlePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) (models.GooglePlacePhotoReference, bool) {
		if googlePlacePhotoReferenceEntity == nil {
			return models.GooglePlacePhotoReference{}, false
		}

		return NewGooglePlacePhotoReferenceFromEntity(*googlePlacePhotoReferenceEntity, googlePlacePhotoAttributionSlice), true
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
