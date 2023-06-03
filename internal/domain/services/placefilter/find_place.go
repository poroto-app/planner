package placefilter

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

func (f PlacesFilter) Find(findFunc func(place places.Place) bool) *places.Place {
	for _, place := range f.placesToFilter {
		if findFunc(place) {
			copy := place
			return &copy
		}
	}
	return nil
}
