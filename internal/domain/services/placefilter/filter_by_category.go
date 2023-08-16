package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FilterByCategory カテゴリに基づいて場所をフィルタリングする。
// includeGivenCategories がtrueの場合は、指定されたカテゴリに含まれる場所のみを残す。
func FilterByCategory(placesToFilter []places.Place, categories []models.LocationCategory, includeGivenCategories bool) []places.Place {
	var subCategories []string
	for _, category := range categories {
		subCategories = append(subCategories, category.SubCategories...)
	}

	return FilterPlaces(placesToFilter, func(place places.Place) bool {
		for _, category := range place.Types {
			if array.IsContain(subCategories, category) {
				return includeGivenCategories
			}
		}
		return !includeGivenCategories
	})
}
