package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
)

// FilterIgnoreCategory ignore categoryを除外する
func FilterIgnoreCategory(placesToFilter []models.PlaceInPlanCandidate) []models.PlaceInPlanCandidate {
	return FilterPlaces(placesToFilter, func(place models.PlaceInPlanCandidate) bool {
		for _, placeType := range place.Google.Types {
			if array.IsContain(models.CategoryIgnore.SubCategories, placeType) {
				return false
			}
		}
		return true
	})
}
