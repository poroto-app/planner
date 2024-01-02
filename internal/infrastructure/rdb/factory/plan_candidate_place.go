package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewPlanCandidatePlaceSliceFromDomainModel(places []models.Place, planCandidateSetId string, planCandidateId string) entities.PlanCandidatePlaceSlice {
	planCandidatePlaces := make(entities.PlanCandidatePlaceSlice, 0, len(places))
	for i, place := range places {
		planCandidatePlaceEntity := NewPlanCandidatePlaceEntityFromDomainModel(place, i, planCandidateSetId, planCandidateId)
		planCandidatePlaces = append(planCandidatePlaces, &planCandidatePlaceEntity)
	}
	return planCandidatePlaces
}

func NewPlanCandidatePlaceEntityFromDomainModel(place models.Place, order int, planCandidateSetId string, planCandidateId string) entities.PlanCandidatePlace {
	return entities.PlanCandidatePlace{
		ID:                 uuid.New().String(),
		PlanCandidateSetID: planCandidateSetId,
		PlanCandidateID:    planCandidateId,
		PlaceID:            place.Id,
		SortOrder:          order,
	}
}
