package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateEntity struct {
	Id        string                  `firestore:"id"`
	Plans     []PlanInCandidateEntity `firestore:"plans"`
	ExpiresAt time.Time               `firestore:"expires_at"`
}

func ToPlanCandidateEntity(planCandidate models.PlanCandidate) PlanCandidateEntity {
	plans := make([]PlanInCandidateEntity, len(planCandidate.Plans))
	for i, plan := range planCandidate.Plans {
		plans[i] = ToPlanInCandidateEntity(plan)
	}

	return PlanCandidateEntity{
		Id:        planCandidate.Id,
		Plans:     plans,
		ExpiresAt: planCandidate.ExpiresAt,
	}
}

func FromPlanCandidateEntity(entity PlanCandidateEntity, metaData PlanCandidateMetaDataV1Entity) models.PlanCandidate {
	plans := make([]models.Plan, len(entity.Plans))
	for i, plan := range entity.Plans {
		plans[i] = fromPlanInCandidateEntity(
			plan.Id,
			plan.Name,
			plan.Places,
			plan.TimeInMinutes,
			plan.Transitions,
		)
	}

	return models.PlanCandidate{
		Id:        entity.Id,
		Plans:     plans,
		MetaData:  FromPlanCandidateMetaDataV1Entity(metaData),
		ExpiresAt: entity.ExpiresAt,
	}
}
