package entity

import (
	"log"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateEntity struct {
	Id                            string                  `firestore:"id"`
	Plans                         []PlanInCandidateEntity `firestore:"plans"`
	CreatedBasedOnCurrentLocation bool                    `firestore:"created_based_on_current_location"`
	ExpiresAt                     time.Time               `firestore:"expires_at"`
}

func ToPlanCandidateEntity(planCandidate models.PlanCandidate) PlanCandidateEntity {
	plans := make([]PlanInCandidateEntity, len(planCandidate.Plans))
	for i, plan := range planCandidate.Plans {
		plans[i] = ToPlanInCandidateEntity(plan)
	}

	return PlanCandidateEntity{
		Id:                            planCandidate.Id,
		Plans:                         plans,
		CreatedBasedOnCurrentLocation: planCandidate.CreatedBasedOnCurrentLocation,
		ExpiresAt:                     planCandidate.ExpiresAt,
	}
}

func FromPlanCandidateEntity(entity PlanCandidateEntity) models.PlanCandidate {
	plans := make([]models.Plan, len(entity.Plans))
	var err error
	for i, plan := range entity.Plans {
		plans[i], err = fromPlanInCandidateEntity(
			plan.Id,
			plan.Name,
			plan.Places,
			plan.PlaceIdsOrdered,
			plan.TimeInMinutes,
			plan.Transitions,
		)
		if err != nil {
			log.Printf("error occur while in converting entity to domain model: [%v]", err)
		}
	}

	return models.PlanCandidate{
		Id:                            entity.Id,
		Plans:                         plans,
		CreatedBasedOnCurrentLocation: entity.CreatedBasedOnCurrentLocation,
		ExpiresAt:                     entity.ExpiresAt,
	}
}
