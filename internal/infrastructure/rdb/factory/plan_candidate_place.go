package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlanCandidatePlaceSliceFromDomainModel(places []models.Place, planCandidateSetId string, planCandidateId string) generated.PlanCandidatePlaceSlice {
	planCandidatePlaces := make(generated.PlanCandidatePlaceSlice, 0, len(places))
	for i, place := range places {
		planCandidatePlaceEntity := NewPlanCandidatePlaceEntityFromDomainModel(place, i, planCandidateSetId, planCandidateId)
		planCandidatePlaces = append(planCandidatePlaces, &planCandidatePlaceEntity)
	}
	return planCandidatePlaces
}

func NewPlanCandidatePlaceEntityFromDomainModel(place models.Place, order int, planCandidateSetId string, planCandidateId string) generated.PlanCandidatePlace {
	return generated.PlanCandidatePlace{
		ID:                 uuid.New().String(),
		PlanCandidateSetID: planCandidateSetId,
		PlanCandidateID:    planCandidateId,
		PlaceID:            place.Id,
		SortOrder:          order,
	}
}
