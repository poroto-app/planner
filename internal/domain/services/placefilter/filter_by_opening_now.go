package placefilter

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

func FilterByOpeningNow(placesToFilter []places.Place) []places.Place {
	return FilterPlaces(placesToFilter, func(place places.Place) bool {
		return place.OpenNow
	})
}
