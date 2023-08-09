package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func FilterWithinDistanceRange(
	placesToFilter []places.Place,
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
) []places.Place {
	return FilterPlaces(placesToFilter, func(place places.Place) bool {
		distance := currentLocation.DistanceInMeter(place.Location.ToGeoLocation())
		return startInMeter <= distance && distance < endInMeter
	})
}
