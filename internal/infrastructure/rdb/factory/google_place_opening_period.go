package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"time"
)

func NewGooglePlaceOpeningPeriodFromEntity(googlePlaceOpeningPeriodEntity generated.GooglePlaceOpeningPeriod) models.GooglePlaceOpeningPeriod {
	return models.GooglePlaceOpeningPeriod{
		DayOfWeekOpen:  time.Weekday(googlePlaceOpeningPeriodEntity.OpenDay).String(),
		DayOfWeekClose: time.Weekday(googlePlaceOpeningPeriodEntity.CloseDay).String(),
		OpeningTime:    googlePlaceOpeningPeriodEntity.OpenTime,
		ClosingTime:    googlePlaceOpeningPeriodEntity.CloseTime,
	}
}

func NewGooglePlaceOpeningPeriodEntityFromDomainModel(googlePlaceOpeningPeriod models.GooglePlaceOpeningPeriod, googlePlaceId string) generated.GooglePlaceOpeningPeriod {
	return generated.GooglePlaceOpeningPeriod{
		ID:            uuid.New().String(),
		OpenDay:       int(weekdayFromWeekdayString(googlePlaceOpeningPeriod.DayOfWeekOpen)),
		CloseDay:      int(weekdayFromWeekdayString(googlePlaceOpeningPeriod.DayOfWeekClose)),
		OpenTime:      googlePlaceOpeningPeriod.OpeningTime,
		CloseTime:     googlePlaceOpeningPeriod.ClosingTime,
		GooglePlaceID: googlePlaceId,
	}
}

func NewGooglePlaceOpeningPeriodSliceFromGooglePlaceDetail(googlePlaceDetail models.GooglePlaceDetail, googlePlaceId string) generated.GooglePlaceOpeningPeriodSlice {
	if googlePlaceDetail.OpeningHours == nil || len(googlePlaceDetail.OpeningHours.Periods) == 0 {
		return nil
	}

	var googlePlaceOpeningPeriodEntities generated.GooglePlaceOpeningPeriodSlice
	for _, googlePlaceOpeningPeriod := range googlePlaceDetail.OpeningHours.Periods {
		gpop := NewGooglePlaceOpeningPeriodEntityFromDomainModel(googlePlaceOpeningPeriod, googlePlaceId)
		googlePlaceOpeningPeriodEntities = append(googlePlaceOpeningPeriodEntities, &gpop)
	}
	return googlePlaceOpeningPeriodEntities
}

func NewGooglePlaceOpeningPeriodSliceFromGooglePlace(googlePlace models.GooglePlace) generated.GooglePlaceOpeningPeriodSlice {
	if googlePlace.PlaceDetail == nil {
		return nil
	}
	return NewGooglePlaceOpeningPeriodSliceFromGooglePlaceDetail(*googlePlace.PlaceDetail, googlePlace.PlaceId)
}

// TODO: OpeningPeriod に time.Weekday 型で持たせる
func weekdayFromWeekdayString(weekdayString string) time.Weekday {
	switch weekdayString {
	case "Sunday":
		return time.Sunday
	case "Monday":
		return time.Monday
	case "Tuesday":
		return time.Tuesday
	case "Wednesday":
		return time.Wednesday
	case "Thursday":
		return time.Thursday
	case "Friday":
		return time.Friday
	case "Saturday":
		return time.Saturday
	default:
		panic("invalid weekday string")
	}
}
