package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlaceDetailFromPlaceDetailEntity(placeDetail places.PlaceDetail) models.GooglePlaceDetail {
	reviews := GooglePlaceReviewsFromPlaceDetail(placeDetail)

	openingPeriods := GooglePlaceOpeningPeriodsFromPlaceDetail(placeDetail)

	return models.GooglePlaceDetail{
		Reviews:      reviews,
		OpeningHours: &openingPeriods,
	}
}
