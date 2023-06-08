package placefilter

import (
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type PlacesFilter struct {
	placesToFilter []places.Place
}

func NewPlacesFilter(placesToFilter []places.Place) PlacesFilter {
	return PlacesFilter{
		placesToFilter: placesToFilter,
	}
}

func (f PlacesFilter) Places() []places.Place {
	return f.placesToFilter
}

func (f PlacesFilter) FilterPlaces(filterFunc func(place places.Place) bool) PlacesFilter {
	var values []places.Place
	for _, place := range f.placesToFilter {
		if filterFunc(place) {
			values = append(values, place)
		}
	}
	return PlacesFilter{
		placesToFilter: values,
	}
}
