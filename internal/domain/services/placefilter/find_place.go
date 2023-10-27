package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func Find(placesToSearch []models.PlaceInPlanCandidate, findFunc func(place models.PlaceInPlanCandidate) bool) *models.PlaceInPlanCandidate {
	for _, place := range placesToSearch {
		if findFunc(place) {
			copy := place
			return &copy
		}
	}
	return nil
}

func FindById(placesToSearch []models.PlaceInPlanCandidate, placeId string) *models.PlaceInPlanCandidate {
	return Find(placesToSearch, func(place models.PlaceInPlanCandidate) bool {
		return place.Id == placeId
	})
}
