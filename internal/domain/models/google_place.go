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
	PhotoReferences  []string
	OpenNow          bool
	PriceLevel       int
	Rating           float32
	UserRatingsTotal int
	Photos           *[]GooglePlacePhoto
	PlaceDetail      *GooglePlaceDetail
}

func (g GooglePlace) Images() []Image {
	if g.Photos == nil {
		return nil
	}

	var images []Image
	for _, photo := range *g.Photos {
		image, err := NewImage(photo.Small, photo.Large)
		if err != nil {
			continue
		}
		images = append(images, *image)
	}

	return images
}

func (g GooglePlace) ToPlaceInPlanCandidate(placeId string) PlaceInPlanCandidate {
	return PlaceInPlanCandidate{
		Id:     placeId,
		Google: g,
	}
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
	if g.PlaceDetail == nil || g.PlaceDetail.OpeningHours == nil {
		return false, fmt.Errorf("opening hours is not set")
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
