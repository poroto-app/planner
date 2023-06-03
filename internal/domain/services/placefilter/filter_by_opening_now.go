package placefilter

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

func (f PlacesFilter) FilterByOpeningNow() PlacesFilter {
	return f.FilterPlaces(func(place places.Place) bool {
		return place.OpenNow
	})
}
