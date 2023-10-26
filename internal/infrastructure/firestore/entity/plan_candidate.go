package entity

import (
	"log"
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

func FromPlanCandidateEntity(entity PlanCandidateEntity, metaData PlanCandidateMetaDataV1Entity, places []models.PlaceInPlanCandidate) models.PlanCandidate {
	plans := make([]models.Plan, 0)
	for _, planEntity := range entity.Plans {
		plan, err := FromPlanInCandidateEntity(
			planEntity.Id,
			planEntity.Name,
			places,
			planEntity.PlaceIdsOrdered,
			planEntity.TimeInMinutes,
		)
		if err != nil {
			log.Printf("error occur while in converting entity.PlanCandidateEntity to models.PlanCandidate: %v", err)
			continue
		}

		// エラーを含むプランが存在した場合，正常なプランだけを返す
		plans = append(plans, *plan)
	}

	return models.PlanCandidate{
		Id:        entity.Id,
		Plans:     plans,
		MetaData:  FromPlanCandidateMetaDataV1Entity(metaData),
		ExpiresAt: entity.ExpiresAt,
	}
}
