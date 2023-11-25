package entity

import "poroto.app/poroto/planner/internal/domain/models"

type GooglePlaceOpeningHoursEntity struct {
	OpeningHoursPeriods []GooglePlaceOpeningPeriodEntity `firestore:"periods"`
}

type GooglePlaceOpeningPeriodEntity struct {
	DayOfWeekOpen  string `firestore:"open_day"`
	DayOfWeekClose string `firestore:"close_day"`
	TimeOpen       string `firestore:"open_time"`
	TimeClose      string `firestore:"close_time"`
}

func GooglePlaceOpeningsEntityFromGooglePlaceOpeningHours(openingHours models.GooglePlaceOpeningHours) GooglePlaceOpeningHoursEntity {
	var entities []GooglePlaceOpeningPeriodEntity
	for _, period := range openingHours.Periods {
		entities = append(entities, GooglePlaceOpeningPeriodEntity{
			DayOfWeekOpen:  period.DayOfWeekOpen,
			DayOfWeekClose: period.DayOfWeekClose,
			TimeOpen:       period.OpeningTime,
			TimeClose:      period.ClosingTime,
		})
	}

	return GooglePlaceOpeningHoursEntity{
		OpeningHoursPeriods: entities,
	}
}

func (g GooglePlaceOpeningHoursEntity) ToGooglePlaceOpeningHours() models.GooglePlaceOpeningHours {
	var periods []models.GooglePlaceOpeningPeriod
	for _, period := range g.OpeningHoursPeriods {
		periods = append(periods, models.GooglePlaceOpeningPeriod{
			DayOfWeekOpen:  period.DayOfWeekOpen,
			DayOfWeekClose: period.DayOfWeekClose,
			OpeningTime:    period.TimeOpen,
			ClosingTime:    period.TimeClose,
		})
	}

	return models.GooglePlaceOpeningHours{
		Periods: periods,
	}
}
