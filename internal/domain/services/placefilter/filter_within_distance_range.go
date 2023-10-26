package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterWithinDistanceRange(
	placesToFilter []models.PlaceInPlanCandidate,
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
) []models.PlaceInPlanCandidate {
	return FilterPlaces(placesToFilter, func(place models.PlaceInPlanCandidate) bool {
		distance := currentLocation.DistanceInMeter(place.Location())
		return startInMeter <= distance && distance < endInMeter
	})
}
