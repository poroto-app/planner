package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlaceDetailFromPlaceDetailEntity(placeDetail places.PlaceDetail) models.GooglePlaceDetail {
	openingPeriods := GooglePlaceOpeningPeriodsFromPlaceDetail(placeDetail)

	return models.GooglePlaceDetail{
		OpeningHours: &openingPeriods,
	}
}
