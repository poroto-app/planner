package models

import (
	"fmt"
)

// Place 場所の情報
type Place struct {
	Id                 string               `json:"id"`
	GooglePlaceId      *string              `json:"google_place_id"`
	Name               string               `json:"name"`
	Location           GeoLocation          `json:"location"`
	Images             []Image              `json:"images"`
	Categories         []LocationCategory   `json:"categories"`
	GooglePlaceReviews *[]GooglePlaceReview `json:"google_place_reviews"`
	PriceLevel         int                  `json:"price_level"`
}

const (
	maxPriceOfLevel1 = 1000
	maxPriceOfLevel2 = 3000
	maxPriceOfLevel3 = 10000
	maxPriceOfLevel4 = 30000
)

func (p Place) MainCategory() *LocationCategory {
	if len(p.Categories) == 0 {
		return nil
	}
	return &p.Categories[0]
}

func (p Place) EstimatedStayDuration() uint {
	categoryMain := p.MainCategory()
	if categoryMain == nil {
		return 0
	}
	return categoryMain.EstimatedStayDuration
}

func (p Place) EstimatedPriceRange() (priceRangeMin, priceRangeMax int, err error) {
	switch p.PriceLevel {
	case 0:
		return 0, 0, nil

	case 1:
		return 0, maxPriceOfLevel1, nil
	case 2:
		return maxPriceOfLevel1, maxPriceOfLevel2, nil
	case 3:
		return maxPriceOfLevel2, maxPriceOfLevel3, nil
	case 4:
		return maxPriceOfLevel3, maxPriceOfLevel4, nil
	default:
		return 0, 0, fmt.Errorf("invalid price level: %d", p.PriceLevel)
	}
}
