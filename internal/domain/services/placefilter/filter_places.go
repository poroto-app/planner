package placefilter

import (
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func FilterPlaces(placesToFilter []places.Place, filterFunc func(place places.Place) bool) []places.Place {
	values := make([]places.Place, 0)
	for _, place := range placesToFilter {
		if filterFunc(place) {
			values = append(values, place)
		}
	}
	return values
}
