package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
)

func FilterDuplicated(placesToFilter []models.GooglePlace) []models.GooglePlace {
	var placeIdsInResult []string
	var placesFiltered []models.GooglePlace
	for _, place := range placesToFilter {
		if !array.IsContain(placeIdsInResult, place.PlaceId) {
			placeIdsInResult = append(placeIdsInResult, place.PlaceId)
			placesFiltered = append(placesFiltered, place)
		}
	}
	return placesFiltered
}
