package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlaceOpeningPeriodsFromPlaceDetail(placeDetail places.PlaceDetail) []models.GooglePlaceOpeningPeriod {
	var openingPeriods []models.GooglePlaceOpeningPeriod
	for _, period := range placeDetail.OpeningHours.Periods {
		openingPeriods = append(openingPeriods, models.GooglePlaceOpeningPeriod{
			DayOfWeekOpen:  period.Open.Day.String(),
			DayOfWeekClose: period.Close.Day.String(),
			OpeningTime:    period.Open.Time,
			ClosingTime:    period.Close.Time,
		})
	}
	return openingPeriods
}
