package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanInCandidateEntity PlanCandidateEntityに含まれるPlan
// MEMO: PlanEntityを用いると、CreatedAtとUpdatedAtが含まれてしまうため、別の構造体を利用している
type PlanInCandidateEntity struct {
	Id     string        `firestore:"id"`
	Name   string        `firestore:"name"`
	Places []PlaceEntity `firestore:"places"`
	// MEMO: Firestoreではuintをサポートしていないため，intにしている
	TimeInMinutes int `firestore:"time_in_minutes"`
}

func ToPlanInCandidateEntity(plan models.Plan) PlanInCandidateEntity {
	ps := make([]PlaceEntity, len(plan.Places))
	for i, place := range plan.Places {
		ps[i] = ToPlaceEntity(place)
	}

	return PlanInCandidateEntity{
		Id:            plan.Id,
		Name:          plan.Name,
		Places:        ps,
		TimeInMinutes: int(plan.TimeInMinutes),
	}
}

func fromPlanInCandidateEntity(
	id string,
	name string,
	places []PlaceEntity,
	timeInMinutes int,
) models.Plan {
	return fromPlanEntity(
		id,
		name,
		places,
		timeInMinutes,
	)
}
