package models

type GooglePlaceDetail struct {
	OpeningHours    *GooglePlaceOpeningHours
	Reviews         []GooglePlaceReview
	PhotoReferences []GooglePlacePhotoReference
}
