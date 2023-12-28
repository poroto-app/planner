package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"time"
)

func NewGooglePlaceOpeningPeriodFromEntity(googlePlaceOpeningPeriodEntity entities.GooglePlaceOpeningPeriod) models.GooglePlaceOpeningPeriod {
	return models.GooglePlaceOpeningPeriod{
		DayOfWeekOpen:  time.Weekday(googlePlaceOpeningPeriodEntity.OpenDay).String(),
		DayOfWeekClose: time.Weekday(googlePlaceOpeningPeriodEntity.CloseDay).String(),
		OpeningTime:    googlePlaceOpeningPeriodEntity.OpenTime,
		ClosingTime:    googlePlaceOpeningPeriodEntity.CloseTime,
	}
}

func NewGooglePlaceOpeningPeriodEntityFromDomainModel(googlePlaceOpeningPeriod models.GooglePlaceOpeningPeriod, googlePlaceId string) (*entities.GooglePlaceOpeningPeriod, error) {
	return &entities.GooglePlaceOpeningPeriod{
		ID:            uuid.New().String(),
		OpenDay:       int(weekdayFromWeekdayString(googlePlaceOpeningPeriod.DayOfWeekOpen)),
		CloseDay:      int(weekdayFromWeekdayString(googlePlaceOpeningPeriod.DayOfWeekClose)),
		OpenTime:      googlePlaceOpeningPeriod.OpeningTime,
		CloseTime:     googlePlaceOpeningPeriod.ClosingTime,
		GooglePlaceID: googlePlaceId,
	}, nil
}

func NewGooglePlaceOpeningPeriodSliceFromGooglePlaceDetail(googlePlaceDetail models.GooglePlaceDetail, googlePlaceId string) (entities.GooglePlaceOpeningPeriodSlice, error) {
	if googlePlaceDetail.OpeningHours == nil || len(googlePlaceDetail.OpeningHours.Periods) == 0 {
		return nil, nil
	}

	var googlePlaceOpeningPeriodEntities entities.GooglePlaceOpeningPeriodSlice
	for _, googlePlaceOpeningPeriod := range googlePlaceDetail.OpeningHours.Periods {
		gpop, err := NewGooglePlaceOpeningPeriodEntityFromDomainModel(googlePlaceOpeningPeriod, googlePlaceId)
		if err != nil {
			return nil, err
		}
		googlePlaceOpeningPeriodEntities = append(googlePlaceOpeningPeriodEntities, gpop)
	}
	return googlePlaceOpeningPeriodEntities, nil
}

func NewGooglePlaceOpeningPeriodSliceFromGooglePlace(googlePlace models.GooglePlace) (entities.GooglePlaceOpeningPeriodSlice, error) {
	if googlePlace.PlaceDetail == nil {
		return nil, nil
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
