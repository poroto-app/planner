package places

import (
	"googlemaps.github.io/maps"
)

type Place struct {
	PlaceID          string   `firestore:"place_id"`
	Name             string   `firestore:"name"`
	Types            []string `firestore:"types"`
	Location         Location `firestore:"location"`
	PhotoReferences  []string `firestore:"photo_references"`
	OpenNow          bool     `firestore:"open_now"`
	Rating           float32  `firestore:"rating"`
	UserRatingsTotal int      `firestore:"user_ratings_total"`
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
) Place {
	return Place{
		PlaceID:          placeID,
		Name:             name,
		Types:            types,
		PhotoReferences:  photoReferences,
		OpenNow:          openNow,
		Rating:           rating,
		UserRatingsTotal: userRatingsTotal,
		Location: Location{
			Latitude:  geometry.Location.Lat,
			Longitude: geometry.Location.Lng,
		},
	}
}
