package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateEntity struct {
	Id        string       `firestore:"id"`
	Plans     []PlanEntity `firestore:"plans"`
	ExpiresAt time.Time    `firestore:"expires_at"`
}

func ToPlanCandidateEntity(planCandidate models.PlanCandidate) PlanCandidateEntity {
	plans := make([]PlanEntity, len(planCandidate.Plans))
	for i, plan := range planCandidate.Plans {
		plans[i] = ToPlanEntity(plan)
	}

	return PlanCandidateEntity{
		Id:        planCandidate.Id,
		Plans:     plans,
		ExpiresAt: planCandidate.ExpiresAt,
	}
}

func FromPlanCandidateEntity(entity PlanCandidateEntity) models.PlanCandidate {
	plans := make([]models.Plan, len(entity.Plans))
	for i, plan := range entity.Plans {
		plans[i] = FromPlanEntity(plan)
	}

	return models.PlanCandidate{
		Id:        entity.Id,
		Plans:     plans,
		ExpiresAt: entity.ExpiresAt,
	}
}
