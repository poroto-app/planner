package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewPlanCandidatePlaceSliceFromDomainModel(places []models.Place, planCandidateSetId string, planCandidateId string) (planCandidatePlaceSlice entities.PlanCandidatePlaceSlice) {
	planCandidatePlaces := make(entities.PlanCandidatePlaceSlice, 0, len(places))
	for i, place := range places {
		planCandidatePlaces = append(planCandidatePlaces, &entities.PlanCandidatePlace{
			ID:                 uuid.New().String(),
			PlanCandidateSetID: planCandidateSetId,
			PlanCandidateID:    planCandidateId,
			PlaceID:            place.Id,
			Order:              i,
		})
	}
	return planCandidatePlaces
}
