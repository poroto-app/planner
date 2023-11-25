package places

import (
	"googlemaps.github.io/maps"
)

type Place struct {
	PlaceID          string
	Name             string
	Types            []string
	Location         Location
	PhotoReferences  []string
	OpenNow          bool
	Rating           float32
	UserRatingsTotal int
	PriceLevel       int
	PlaceDetail      *PlaceDetail
}

type PlaceDetail struct {
	Reviews      []maps.PlaceReview
	Photos       []maps.Photo
	OpeningHours *maps.OpeningHours
}

type Location struct {
	Latitude  float64 `firestore:"latitude"`
	Longitude float64 `firestore:"longitude"`
}

func createPlace(
	placeID string,
	name string,
	types []string,
	geometry maps.AddressGeometry,
	photoReferences []string,
	openNow bool,
	rating float32,
	userRatingsTotal int,
	priceLevel int,
) Place {
	return Place{
		PlaceID:          placeID,
		Name:             name,
		Types:            types,
		PhotoReferences:  photoReferences,
		OpenNow:          openNow,
		Rating:           rating,
		UserRatingsTotal: userRatingsTotal,
		PriceLevel:       priceLevel,
		Location: Location{
			Latitude:  geometry.Location.Lat,
			Longitude: geometry.Location.Lng,
		},
	}
}

func createPlaceDetail(
	reviews []maps.PlaceReview,
	photos []maps.Photo,
	openingHours *maps.OpeningHours,
) PlaceDetail {
	return PlaceDetail{
		Reviews:      reviews,
		Photos:       photos,
		OpeningHours: openingHours,
	}
}
