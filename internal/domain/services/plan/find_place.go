package plan

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

func (s PlanService) findPlace(
	placesToScan []places.Place,
	findFunc func(place places.Place) bool,
) *places.Place {
	for _, place := range placesToScan {
		if findFunc(place) {
			copy := place
			return &copy
		}
	}

	return nil
}
