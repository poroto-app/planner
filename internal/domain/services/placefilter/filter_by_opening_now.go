package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterByOpeningNow(placesToFilter []models.PlaceInPlanCandidate) []models.PlaceInPlanCandidate {
	return FilterPlaces(placesToFilter, func(place models.PlaceInPlanCandidate) bool {
		return place.Google.OpenNow
	})
}
