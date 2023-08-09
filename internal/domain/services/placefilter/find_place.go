package placefilter

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

func Find(placesToSearch []places.Place, findFunc func(place places.Place) bool) *places.Place {
	for _, place := range placesToSearch {
		if findFunc(place) {
			copy := place
			return &copy
		}
	}
	return nil
}

func FindById(placesToSearch []places.Place, placeId string) *places.Place {
	return Find(placesToSearch, func(place places.Place) bool {
		return place.PlaceID == placeId
	})
}
