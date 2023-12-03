package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterByOpeningNow(placesToFilter []models.Place) []models.Place {
	return FilterPlaces(placesToFilter, func(place models.Place) bool {
		return place.Google.OpenNow
	})
}
