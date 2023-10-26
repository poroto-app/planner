package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterByRating(placesToFilter []models.PlaceInPlanCandidate, lowestRating float32, lowestUserRatingsTotal int) []models.PlaceInPlanCandidate {
	return FilterPlaces(placesToFilter, func(place models.PlaceInPlanCandidate) bool {
		return place.Google.Rating >= lowestRating && place.Google.UserRatingsTotal >= lowestUserRatingsTotal
	})
}
