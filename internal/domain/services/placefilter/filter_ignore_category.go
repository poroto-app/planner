package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FilterIgnoreCategory ignore categoryを除外する
func (f PlacesFilter) FilterIgnoreCategory() PlacesFilter {
	return f.FilterPlaces(func(place places.Place) bool {
		for _, category := range place.Types {
			if array.IsContain(models.CategoryIgnore.SubCategories, category) {
				return false
			}
		}
		return true
	})
}
