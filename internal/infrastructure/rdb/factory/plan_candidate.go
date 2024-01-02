package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func PlanCandidateSliceFromDomainModel(planCandidates []models.Plan, planCandidateSetId string) entities.PlanCandidateSlice {
	var planCandidateSlice entities.PlanCandidateSlice
	for _, planCandidate := range planCandidates {
		planCandidateEntity := PlanCandidateEntityFromDomainModel(planCandidate, planCandidateSetId)
		planCandidateSlice = append(planCandidateSlice, &planCandidateEntity)
	}
	return planCandidateSlice
}

func PlanCandidateEntityFromDomainModel(planCandidate models.Plan, planCandidateSetId string) entities.PlanCandidate {
	return entities.PlanCandidate{
		ID:                 planCandidate.Id,
		PlanCandidateSetID: planCandidateSetId,
		Name:               planCandidate.Name,
	}
}
