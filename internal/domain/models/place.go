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
