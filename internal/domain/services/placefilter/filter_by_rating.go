package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterByRating(placesToFilter []models.Place, lowestRating float32, lowestUserRatingsTotal int) []models.Place {
	return FilterPlaces(placesToFilter, func(place models.Place) bool {
		return place.Google.Rating >= lowestRating && place.Google.UserRatingsTotal >= lowestUserRatingsTotal
	})
}
