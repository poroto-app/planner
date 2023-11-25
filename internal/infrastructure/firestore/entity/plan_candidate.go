package entity

import (
	"log"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateEntity struct {
	Id        string    `firestore:"id"`
	PlanIds   []string  `firestore:"plan_ids"`
	ExpiresAt time.Time `firestore:"expires_at"`
}

func ToPlanCandidateEntity(planCandidate models.PlanCandidate) PlanCandidateEntity {
	plansIds := make([]string, len(planCandidate.Plans))
	for i, plan := range planCandidate.Plans {
		plansIds[i] = plan.Id
	}

	return PlanCandidateEntity{
		Id:        planCandidate.Id,
		PlanIds:   plansIds,
		ExpiresAt: planCandidate.ExpiresAt,
	}
}

func FromPlanCandidateEntity(entity PlanCandidateEntity, metaData PlanCandidateMetaDataV1Entity, planEntities []PlanInCandidateEntity, places []models.PlaceInPlanCandidate) models.PlanCandidate {
	var plans []models.Plan
	for _, planId := range entity.PlanIds {
		for _, place := range planEntities {
			if place.Id != planId {
				continue
			}

			plan, err := FromPlanInCandidateEntity(planId, place.Name, places, place.PlaceIdsOrdered)
			if err != nil {
				log.Printf("error while converting entity.PlanCandidateEntity to models.PlanCandidate: %v", err)
			}

			plans = append(plans, *plan)
		}
	}

	return models.PlanCandidate{
		Id:        entity.Id,
		Plans:     plans,
		MetaData:  FromPlanCandidateMetaDataV1Entity(metaData),
		ExpiresAt: entity.ExpiresAt,
	}
}
