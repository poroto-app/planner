package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterDuplicated(placesToFilter []models.Place) []models.Place {
	var placeIdsInResult []string
	var placesFiltered []models.Place
	for _, place := range placesToFilter {
		if !array.IsContain(placeIdsInResult, place.Id) {
			placeIdsInResult = append(placeIdsInResult, place.Id)
			placesFiltered = append(placesFiltered, place)
		}
	}
	return placesFiltered
}
