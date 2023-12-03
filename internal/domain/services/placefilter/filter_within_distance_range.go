package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterWithinDistanceRange(
	placesToFilter []models.Place,
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
) []models.Place {
	return FilterPlaces(placesToFilter, func(place models.Place) bool {
		distance := currentLocation.DistanceInMeter(place.Location)
		return startInMeter <= distance && distance < endInMeter
	})
}
