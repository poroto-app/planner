package placefilter

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

func (f PlacesFilter) FilterByOpeningNow() PlacesFilter {
	f.placesToFilter = f.filterPlaces(func(place places.Place) bool {
		return place.OpenNow
	})
	return f
}
