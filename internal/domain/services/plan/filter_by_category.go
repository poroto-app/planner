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

	var placesInCategory []places.Place
	for _, place := range placesToFilter {
		if array.HasIntersection(place.Types, categoriesSlice) {
			placesInCategory = append(placesInCategory, place)
		}
	}

	return placesInCategory
}
