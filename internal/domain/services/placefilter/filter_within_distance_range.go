package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (f PlacesFilter) FilterWithinDistanceRange(
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
) PlacesFilter {
	f.placesToFilter = f.filterPlaces(func(place places.Place) bool {
		distance := currentLocation.DistanceInMeter(place.Location.ToGeoLocation())
		return startInMeter <= distance && distance < endInMeter
	})
	return f
}
