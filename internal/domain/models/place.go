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

// EstimatedPriceRange 価格帯を推定する
// SEE: https://developers.google.com/maps/documentation/places/web-service/details?hl=ja#Place-price_level
func (p Place) EstimatedPriceRange() (priceRange *PriceRange) {
	switch p.PriceLevel {
	case 0:
		// TODO: 飲食店でprice_levelが0の場合は、価格帯が不明なので、nilを返す
		return &PriceRange{
			Min:              0,
			Max:              0,
			GooglePriceLevel: p.PriceLevel,
		}
	case 1:
		return &PriceRange{
			Min:              0,
			Max:              maxPriceOfLevel1,
			GooglePriceLevel: p.PriceLevel,
		}
	case 2:
		return &PriceRange{
			Min:              0,
			Max:              maxPriceOfLevel2,
			GooglePriceLevel: p.PriceLevel,
		}
	case 3:
		return &PriceRange{
			Min:              0,
			Max:              maxPriceOfLevel3,
			GooglePriceLevel: p.PriceLevel,
		}
	case 4:
		return &PriceRange{
			Min:              0,
			Max:              maxPriceOfLevel4,
			GooglePriceLevel: p.PriceLevel,
		}
	default:
		return nil
	}
}
