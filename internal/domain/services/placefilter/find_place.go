package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

func Find(placesToSearch []models.GooglePlace, findFunc func(place models.GooglePlace) bool) *models.GooglePlace {
	for _, place := range placesToSearch {
		if findFunc(place) {
			copy := place
			return &copy
		}
	}
	return nil
}

func FindById(placesToSearch []models.GooglePlace, placeId string) *models.GooglePlace {
	return Find(placesToSearch, func(place models.GooglePlace) bool {
		return place.PlaceId == placeId
	})
}
