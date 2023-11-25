package models

// Place 場所の情報
type Place struct {
	Id                 string               `json:"id"`
	Google             GooglePlace          `json:"google"`
	GooglePlaceId      *string              `json:"google_place_id"`
	Name               string               `json:"name"`
	Location           GeoLocation          `json:"location"`
	Images             []Image              `json:"images"`
	Categories         []LocationCategory   `json:"categories"`
	GooglePlaceReviews *[]GooglePlaceReview `json:"google_place_reviews"`
	PriceLevel         int                  `json:"price_level"`
}

func NewPlaceFromGooglePlace(placeId string, googlePlace GooglePlace) Place {
	return Place{
		Id:       placeId,
		Google:   googlePlace,
		Name:     googlePlace.Name,
		Location: googlePlace.Location,
	}
}

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
func (p Place) EstimatedPriceRange() (priceRange *PriceRange) {
	// TODO: 飲食店でprice_levelが0の場合は、価格帯が不明なので、nilを返す
	return PriceRangeFromGooglePriceLevel(p.PriceLevel)
}
