package models

import (
	"fmt"
	"time"
)

type GooglePlace struct {
	PlaceId          string
	Name             string
	Types            []string
	Location         GeoLocation
	PhotoReferences  []GooglePlacePhotoReference
	PriceLevel       int
	Rating           float32
	UserRatingsTotal int
	Vicinity         *string
	FormattedAddress *string
	Photos           *[]GooglePlacePhoto
	PlaceDetail      *GooglePlaceDetail
}

// IndexOfCategory は Types 中の `category` に対応する Type のインデックスを返す
func (g GooglePlace) IndexOfCategory(category LocationCategory) int {
	for i, placeType := range g.Types {
		c := CategoryOfSubCategory(placeType)
		if c.Name == category.Name {
			return i
		}
	}
	return -1
}

func (g GooglePlace) IsOpening(at time.Time) (bool, error) {
	if g.PlaceDetail == nil {
		return false, fmt.Errorf("place detail is not fetched")
	}

	// OpeningHoursがない場合は開いているとみなす
	// (営業時間が不明な場所がスキップされてしまうとプランに含まれなくなるため)
	if g.PlaceDetail.OpeningHours == nil {
		return true, nil
	}

	for _, openingPeriod := range g.PlaceDetail.OpeningHours.Periods {
		weekday := at.Weekday()
		isOpeningPeriodOfToday := openingPeriod.DayOfWeekOpen == weekday.String()
		if !isOpeningPeriodOfToday {
			continue
		}

		openingTime, err := openingPeriod.OpeningTimeHHMM()
		if err != nil {
			return false, fmt.Errorf("error while parsing opening time")
		}

		closingTime, err := openingPeriod.ClosingTimeHHMM()
		if err != nil {
			return false, fmt.Errorf("error while parsing closing time")
		}

		today := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, at.Location())
		openingTimeAtToday := today.Add(time.Hour*time.Duration(openingTime.Hour) + time.Minute*time.Duration(openingTime.Minute))
		closingTimeAtToday := today.Add(time.Hour*time.Duration(closingTime.Hour) + time.Minute*time.Duration(closingTime.Minute))

		if at.After(openingTimeAtToday) && at.Before(closingTimeAtToday) {
			return true, nil
		}
	}

	return false, nil
}
