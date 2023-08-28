package placefilter

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

func FilterByRating(placesToFilter []places.Place, lowestRating float32, lowestUserRatingsTotal int) []places.Place {
	return FilterPlaces(placesToFilter, func(place places.Place) bool {
		return place.Rating >= lowestRating && place.UserRatingsTotal >= lowestUserRatingsTotal
	})
}
