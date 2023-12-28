package factory

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewGooglePlaceDetailFromGooglePlaceEntity(googlePlaceEntity entities.GooglePlace) (*models.GooglePlaceDetail, error) {
	if googlePlaceEntity.R == nil {
		return nil, fmt.Errorf("googlePlaceEntity.R is nil")
	}

	var googlePlaceReviews []models.GooglePlaceReview
	for _, googlePlaceReviewEntity := range googlePlaceEntity.R.GetGooglePlaceReviews() {
		if googlePlaceReviewEntity == nil {
			continue
		}

		googlePlaceReviews = append(googlePlaceReviews, NewGooglePlaceReviewFromEntity(*googlePlaceReviewEntity))
	}

	var googlePlaceOpeningPeriods []models.GooglePlaceOpeningPeriod
	for _, googlePlaceOpeningPeriodEntity := range googlePlaceEntity.R.GetGooglePlaceOpeningPeriods() {
		if googlePlaceOpeningPeriodEntity == nil {
			continue
		}

		googlePlaceOpeningPeriods = append(googlePlaceOpeningPeriods, NewGooglePlaceOpeningPeriodFromEntity(*googlePlaceOpeningPeriodEntity))
	}

	var googlePlacePhotoReferences []models.GooglePlacePhotoReference
	for _, googlePlacePhotoReferenceEntity := range googlePlaceEntity.R.GetGooglePlacePhotoReferences() {
		if googlePlacePhotoReferenceEntity == nil {
			continue
		}

		gpr := NewGooglePlacePhotoReferenceFromEntity(*googlePlacePhotoReferenceEntity, googlePlaceEntity.R.GetGooglePlacePhotoAttributions())
		googlePlacePhotoReferences = append(googlePlacePhotoReferences, gpr)
	}

	if len(googlePlaceReviews) == 0 && len(googlePlaceOpeningPeriods) == 0 && len(googlePlacePhotoReferences) == 0 {
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
