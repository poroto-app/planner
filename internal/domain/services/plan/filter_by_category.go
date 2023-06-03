package plan

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s PlanService) filterByCategory(
	placesToFilter []places.Place,
	categories []models.LocationCategory,
) []places.Place {
	var categoriesSlice []string
	for _, category := range categories {
		categoriesSlice = append(categoriesSlice, category.SubCategories...)
	}

	return s.filterPlaces(placesToFilter, func(place places.Place) bool {
		for _, category := range place.Types {
			if array.IsContain(categoriesSlice, category) {
				return true
			}
		}
		return false
	})
}
