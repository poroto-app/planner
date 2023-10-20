package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterWithinDistanceRange(
	placesToFilter []models.GooglePlace,
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
) []models.GooglePlace {
	return FilterPlaces(placesToFilter, func(place models.GooglePlace) bool {
		distance := currentLocation.DistanceInMeter(place.Location)
		return startInMeter <= distance && distance < endInMeter
	})
}
