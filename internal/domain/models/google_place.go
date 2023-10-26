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
	PriceLevel       *int
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
