package models

type GooglePlaceDetail struct {
	OpeningHours *[]GooglePlaceOpeningPeriod
	Reviews      []GooglePlaceReview
	Photos       []GooglePlacePhoto
}
