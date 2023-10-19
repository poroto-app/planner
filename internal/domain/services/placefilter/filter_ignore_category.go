package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
)

// FilterIgnoreCategory ignore categoryを除外する
func FilterIgnoreCategory(placesToFilter []models.GooglePlace) []models.GooglePlace {
	return FilterPlaces(placesToFilter, func(place models.GooglePlace) bool {
		for _, category := range place.Types {
			if array.IsContain(models.CategoryIgnore.SubCategories, category) {
				return false
			}
		}
		return true
	})
}
