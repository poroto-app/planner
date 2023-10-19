package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterByRating(placesToFilter []models.GooglePlace, lowestRating float32, lowestUserRatingsTotal int) []models.GooglePlace {
	return FilterPlaces(placesToFilter, func(place models.GooglePlace) bool {
		return place.Rating >= lowestRating && place.UserRatingsTotal >= lowestUserRatingsTotal
	})
}
