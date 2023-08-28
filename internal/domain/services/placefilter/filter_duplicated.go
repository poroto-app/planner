package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func FilterDuplicated(placesToFilter []places.Place) []places.Place {
	var placeIdsInResult []string
	var placesFiltered []places.Place
	for _, place := range placesToFilter {
		if !array.IsContain(placeIdsInResult, place.PlaceID) {
			placeIdsInResult = append(placeIdsInResult, place.PlaceID)
			placesFiltered = append(placesFiltered, place)
		}
	}
	return placesFiltered
}
