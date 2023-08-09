package placefilter

import (
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func FilterPlaces(placesToFilter []places.Place, filterFunc func(place places.Place) bool) []places.Place {
	var values []places.Place
	for _, place := range placesToFilter {
		if filterFunc(place) {
			values = append(values, place)
		}
	}
	return values
}
