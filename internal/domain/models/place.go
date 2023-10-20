package models

// Place 場所の情報
type Place struct {
	Id                 string               `json:"id"`
	GooglePlaceId      *string              `json:"google_place_id"`
	Name               string               `json:"name"`
	Location           GeoLocation          `json:"location"`
	Images             []Image              `json:"images"`
	Categories         []LocationCategory   `json:"categories"`
	GooglePlaceReviews *[]GooglePlaceReview `json:"google_place_reviews"`
	PriceLevel         *int                 `json:"price_level"`
}

const (
	limitOfPriceRangeMin         = 0
	thresholdOfLevel0AndLevel1_2 = 1000
	thresholdOfLevel1_2AndLevel3 = 3000
	thresholdOfLevel3AndLevel4   = 10000
	limitOfPriceRangeMax         = 30000
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

func (p Place) EstimatedPriceRange() (priceRangeMin, priceRangeMax *int) {
	switch *p.PriceLevel {
	case 0:
		return toIntPointer(limitOfPriceRangeMin), toIntPointer(thresholdOfLevel0AndLevel1_2)
	case 1, 2:
		return toIntPointer(thresholdOfLevel0AndLevel1_2), toIntPointer(thresholdOfLevel1_2AndLevel3)
	case 3:
		return toIntPointer(thresholdOfLevel1_2AndLevel3), toIntPointer(thresholdOfLevel3AndLevel4)
	case 4:
		return toIntPointer(thresholdOfLevel3AndLevel4), toIntPointer(limitOfPriceRangeMax)
	}
	return nil, nil
}

func toIntPointer(x int) *int {
	return &x
}
