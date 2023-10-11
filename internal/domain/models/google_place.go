package models

type GooglePlace struct {
	PlaceId          string
	Name             string
	Types            []string
	Location         GeoLocation
	PhotoReferences  []string
	OpenNow          bool
	Rating           float32
	UserRatingsTotal int
	Images           *[]Image
	Reviews          *[]GooglePlaceReview
}
