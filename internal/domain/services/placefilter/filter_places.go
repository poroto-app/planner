package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterPlaces(placesToFilter []models.PlaceInPlanCandidate, filterFunc func(place models.PlaceInPlanCandidate) bool) []models.PlaceInPlanCandidate {
	values := make([]models.PlaceInPlanCandidate, 0)
	for _, place := range placesToFilter {
		if filterFunc(place) {
			values = append(values, place)
		}
	}
	return values
}
