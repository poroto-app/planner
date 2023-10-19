package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
)

// FilterByCategory カテゴリに基づいて場所をフィルタリングする。
// includeGivenCategories がtrueの場合は、指定されたカテゴリに含まれる場所のみを残す。
func FilterByCategory(placesToFilter []models.GooglePlace, categories []models.LocationCategory, includeGivenCategories bool) []models.GooglePlace {
	var subCategories []string
	for _, category := range categories {
		subCategories = append(subCategories, category.SubCategories...)
	}

	return FilterPlaces(placesToFilter, func(place models.GooglePlace) bool {
		for _, category := range place.Types {
			if array.IsContain(subCategories, category) {
				return includeGivenCategories
			}
		}
		return !includeGivenCategories
	})
}
