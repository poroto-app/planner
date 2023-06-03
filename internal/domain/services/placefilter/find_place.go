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

func (f PlacesFilter) FindById(placeId string) *places.Place {
	return f.Find(func(place places.Place) bool {
		return place.PlaceID == placeId
	})
}
