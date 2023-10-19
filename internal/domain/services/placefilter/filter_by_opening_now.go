package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterByOpeningNow(placesToFilter []models.GooglePlace) []models.GooglePlace {
	return FilterPlaces(placesToFilter, func(place models.GooglePlace) bool {
		return place.OpenNow
	})
}
