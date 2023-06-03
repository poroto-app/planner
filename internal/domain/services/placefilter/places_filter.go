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

func (f PlacesFilter) Copy() PlacesFilter {
	return PlacesFilter{
		placesToFilter: f.placesToFilter,
	}
}

func (f PlacesFilter) Places() []places.Place {
	return f.placesToFilter
}

func (f PlacesFilter) filterPlaces(filterFunc func(place places.Place) bool) []places.Place {
	var values []places.Place
	for _, place := range f.placesToFilter {
		if filterFunc(place) {
			values = append(values, place)
		}
	}
	return values
}
