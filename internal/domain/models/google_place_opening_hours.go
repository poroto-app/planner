package models

import (
	"fmt"
	"strconv"
)

type GooglePlaceOpeningPeriod struct {
	DayOfWeekOpen  string
	DayOfWeekClose string
	OpeningTime    string
	ClosingTime    string
}

type TimeHHMM struct {
	Hour   int
	Minute int
}

func (g GooglePlaceOpeningPeriod) OpeningTimeHHMM() (*TimeHHMM, error) {
	return parseTimeString(g.OpeningTime)
}

func (g GooglePlaceOpeningPeriod) ClosingTimeHHMM() (*TimeHHMM, error) {
	return parseTimeString(g.ClosingTime)

}

// parseTimeString 0000 ~ 2359 の形式で与えられる時間をパースする
// See: https://developers.google.com/maps/documentation/places/web-service/details?hl=ja#PlaceOpeningHoursPeriodDetail-time
func parseTimeString(timeStr string) (time *TimeHHMM, err error) {
	hour, err := strconv.Atoi(timeStr[:2])
	if err != nil {
		return nil, fmt.Errorf("error while parsing hour: %v", err)
	}

	minute, err := strconv.Atoi(timeStr[2:])
	if err != nil {
		return nil, fmt.Errorf("error while parsing minute: %v", err)
	}

	return &TimeHHMM{
		Hour:   hour,
		Minute: minute,
	}, nil
}
