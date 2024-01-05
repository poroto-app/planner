package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"sort"
)

func NewGooglePlaceTypesFromEntity(googlePlaceTypeEntitySlice entities.GooglePlaceTypeSlice) []string {
	var googlePlaceTypeEntities []entities.GooglePlaceType
	for _, googlePlaceTypeEntity := range googlePlaceTypeEntitySlice {
		if googlePlaceTypeEntity == nil {
			continue
		}
		googlePlaceTypeEntities = append(googlePlaceTypeEntities, *googlePlaceTypeEntity)
	}

	sort.Slice(googlePlaceTypeEntities, func(i, j int) bool {
		return googlePlaceTypeEntities[i].OrderNum < googlePlaceTypeEntities[j].OrderNum
	})

	var googlePlaceTypes []string
	for _, googlePlaceTypeEntity := range googlePlaceTypeEntities {
		googlePlaceTypes = append(googlePlaceTypes, googlePlaceTypeEntity.Type)
	}

	return googlePlaceTypes
}

func NewGooglePlaceTypeSliceFromGooglePlace(googlePlace models.GooglePlace) entities.GooglePlaceTypeSlice {
	var googlePlaceTypeEntities entities.GooglePlaceTypeSlice
	for i, googlePlaceType := range googlePlace.Types {
		googlePlaceTypeEntities = append(googlePlaceTypeEntities, &entities.GooglePlaceType{
			ID:            uuid.New().String(),
			GooglePlaceID: googlePlaceType,
			Type:          googlePlaceType,
			OrderNum:      i,
		})
	}
	return googlePlaceTypeEntities
}
