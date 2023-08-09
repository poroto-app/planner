package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func FilterByCategory(placesToFilter []places.Place, categories []models.LocationCategory) []places.Place {
	var subCategories []string
	for _, category := range categories {
		subCategories = append(subCategories, category.SubCategories...)
	}

	return FilterPlaces(placesToFilter, func(place places.Place) bool {
		for _, category := range place.Types {
			if array.IsContain(subCategories, category) {
				return true
			}
		}
		return false
	})
}
