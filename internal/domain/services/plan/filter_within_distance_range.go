package plan

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s PlanService) filterWithinDistanceRange(
	placesToFilter []places.Place,
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
) []places.Place {
	var placesWithInDistance []places.Place
	for _, place := range placesToFilter {
		distance := currentLocation.DistanceInMeter(place.Location.ToGeoLocation())
		if startInMeter <= distance && distance < endInMeter {
			placesWithInDistance = append(placesWithInDistance, place)
		}
	}
	return placesWithInDistance
}
