package plangen

import (
	"context"
	"log"
	"strconv"
	"time"

	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// isOpeningWithIn は，指定された場所が指定された時間内に開いているかを判定する
func (s Service) isOpeningWithIn(
	ctx context.Context,
	place places.Place,
	startTime time.Time,
	duration time.Duration,
) bool {
	placeOpeningPeriods, err := s.placesApi.FetchPlaceOpeningPeriods(ctx, place.PlaceID)
	if err != nil {
		log.Printf("error while fetching place opening periods: %v\n", err)
		return false
	}

	for _, placeOpeningPeriod := range placeOpeningPeriods {
		weekday := startTime.Weekday()
		isOpeningPeriodOfToday := placeOpeningPeriod.DayOfWeek == weekday.String()
		if !isOpeningPeriodOfToday {
			continue
		}

		openingHour, openingMinute, err := parseTimeString(placeOpeningPeriod.OpeningTime)
		if err != nil {
			log.Println("error while parsing opening time")
			continue
		}

		closingHour, closingMinute, err := parseTimeString(placeOpeningPeriod.ClosingTime)
		if err != nil {
			log.Println("error while parsing closing time")
		}

		today := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
		openingTime := today.Add(time.Hour*time.Duration(openingHour) + time.Minute*time.Duration(openingMinute))
		closingTime := today.Add(time.Hour*time.Duration(closingHour) + time.Minute*time.Duration(closingMinute))

		timeEndOfPlan := startTime.Add(time.Minute * duration)

		// 開店時刻 < 開始時刻 && 終了時刻 < 閉店時刻 の判断
		if startTime.After(openingTime) && timeEndOfPlan.Before(closingTime) {
			return true
		}
	}

	return false
}

// parseTimeString 0000 ~ 2359 の形式で与えられる時間をパースする
// See: https://developers.google.com/maps/documentation/places/web-service/details?hl=ja#PlaceOpeningHoursPeriodDetail-time
func parseTimeString(timeStr string) (hour, minute int, err error) {
	hour, err = strconv.Atoi(timeStr[:2])
	if err != nil {
		return 0, 0, nil
	}

	minute, err = strconv.Atoi(timeStr[2:])
	if err != nil {
		return 0, 0, nil
	}

	return hour, minute, nil
}
