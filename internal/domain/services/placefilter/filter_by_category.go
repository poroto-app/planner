package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
)

// FilterByCategory カテゴリに基づいて場所をフィルタリングする。
// includeGivenCategories がtrueの場合は、指定されたカテゴリに含まれる場所のみを残す。
func FilterByCategory(placesToFilter []models.PlaceInPlanCandidate, categories []models.LocationCategory, includeGivenCategories bool) []models.PlaceInPlanCandidate {
	return FilterPlaces(placesToFilter, func(place models.PlaceInPlanCandidate) bool {
		for _, c := range categories {
			for _, placeTypes := range place.Google.Types {
				if array.IsContain(c.SubCategories, placeTypes) {
					return includeGivenCategories
				}
			}
		}
		return !includeGivenCategories
	})
}
