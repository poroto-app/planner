package factory

import (
	"fmt"
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"sort"
)

func NewPlanPlaceSliceFromDomainMode(planPlaces []models.Place, planId string) generated.PlanPlaceSlice {
	var planPlaceSlice generated.PlanPlaceSlice
	for i, place := range planPlaces {
		planPlaceSlice = append(planPlaceSlice, &generated.PlanPlace{
			ID:        uuid.New().String(),
			PlaceID:   place.Id,
			PlanID:    planId,
			SortOrder: i,
		})
	}
	return planPlaceSlice
}

func NewPlanPlacesFromEntities(
	planPlaceSlice generated.PlanPlaceSlice,
	places []models.Place,
	planId string,
) (*[]models.Place, error) {
	planPlaceEntities := array.MapAndFilter(planPlaceSlice, func(planPlaceEntity *generated.PlanPlace) (generated.PlanPlace, bool) {
		if planPlaceEntity == nil {
			return generated.PlanPlace{}, false
		}
		if planPlaceEntity.PlanID != planId {
			return generated.PlanPlace{}, false
		}
		return *planPlaceEntity, true
	})

	// SortOrder でソート
	sort.Slice(planPlaceEntities, func(i, j int) bool {
		return planPlaceEntities[i].SortOrder < planPlaceEntities[j].SortOrder
	})

	planPlaces, err := array.MapWithErr(planPlaceEntities, func(planPlaceEntity generated.PlanPlace) (*models.Place, error) {
		place, ok := array.Find(places, func(place models.Place) bool {
			return place.Id == planPlaceEntity.PlaceID
		})
		if !ok {
			return nil, fmt.Errorf("place not found: %s", planPlaceEntity.PlaceID)
		}
		return &place, nil
	})
	if err != nil {
		return nil, err
	}

	return planPlaces, nil
}
