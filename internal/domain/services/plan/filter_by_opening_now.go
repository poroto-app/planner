package plan

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

func (s PlanService) filterByOpeningNow(
	placesToFilter []places.Place,
) []places.Place {
	var placesOpeningNow []places.Place
	for _, place := range placesToFilter {
		if place.OpenNow {
			placesOpeningNow = append(placesOpeningNow, place)
		}
	}
	return placesOpeningNow
}
