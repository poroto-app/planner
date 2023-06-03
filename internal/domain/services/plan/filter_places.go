package plan

import (
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s PlanService) filterPlaces(
	placesToFilter []places.Place,
	filterFunc func(place places.Place) bool,
) []places.Place {
	var filteredPlaces []places.Place
	for _, place := range placesToFilter {
		if filterFunc(place) {
			filteredPlaces = append(filteredPlaces, place)
		}
	}
	return filteredPlaces
}
