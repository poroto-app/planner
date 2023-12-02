package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterPlaces(placesToFilter []models.Place, filterFunc func(place models.Place) bool) []models.Place {
	values := make([]models.Place, 0)
	for _, place := range placesToFilter {
		if filterFunc(place) {
			values = append(values, place)
		}
	}
	return values
}
