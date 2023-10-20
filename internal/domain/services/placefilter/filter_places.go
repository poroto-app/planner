package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterPlaces(placesToFilter []models.GooglePlace, filterFunc func(place models.GooglePlace) bool) []models.GooglePlace {
	values := make([]models.GooglePlace, 0)
	for _, place := range placesToFilter {
		if filterFunc(place) {
			values = append(values, place)
		}
	}
	return values
}
