package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlanPlaceSliceFromDomainMode(planPlaces []models.Place, planId string) generated.PlanPlaceSlice {
	var planPlaceSlice generated.PlanPlaceSlice
	for i, place := range planPlaces {
		planPlaceSlice = append(planPlaceSlice, &generated.PlanPlace{
			ID:        place.Id,
			PlaceID:   place.Id,
			PlanID:    planId,
			SortOrder: i,
		})
	}
	return planPlaceSlice
}
