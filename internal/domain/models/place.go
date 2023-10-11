package models

// Place 場所の情報
type Place struct {
	Id                    string               `json:"id"`
	GooglePlaceId         *string              `json:"google_place_id"`
	Name                  string               `json:"name"`
	Location              GeoLocation          `json:"location"`
	Images                []Image              `json:"images"`
	EstimatedStayDuration uint                 `json:"estimated_stay_duration"`
	Categories            []LocationCategory   `json:"categories"`
	GooglePlaceReviews    *[]GooglePlaceReview `json:"google_place_reviews"`
}
