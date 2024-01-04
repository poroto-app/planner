package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func PlanCandidateEntityFromDomainModel(planCandidate models.Plan, planCandidateSetId string, sortOrder int) entities.PlanCandidate {
	return entities.PlanCandidate{
		ID:                 planCandidate.Id,
		PlanCandidateSetID: planCandidateSetId,
		Name:               planCandidate.Name,
		SortOrder:          sortOrder,
	}
}
