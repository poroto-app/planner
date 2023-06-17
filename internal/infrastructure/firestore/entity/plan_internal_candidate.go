package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanInternalCandidateEntity PlanCandidateEntityに含まれるPlan
// MEMO: PlanEntityを用いると、CreatedAtとUpdatedAtが含まれてしまうため、別の構造体を利用している
type PlanInternalCandidateEntity struct {
	Id     string        `firestore:"id"`
	Name   string        `firestore:"name"`
	Places []PlaceEntity `firestore:"places"`
	// MEMO: Firestoreではuintをサポートしていないため，intにしている
	TimeInMinutes int `firestore:"time_in_minutes"`
}

func toPlanCandidateEntity(
	id string,
	name string,
	places []models.Place,
	timeInMinutes uint,
) PlanInternalCandidateEntity {
	ps := make([]PlaceEntity, len(places))
	for i, place := range places {
		ps[i] = ToPlaceEntity(place)
	}

	return PlanInternalCandidateEntity{
		Id:            id,
		Name:          name,
		Places:        ps,
		TimeInMinutes: int(timeInMinutes),
	}
}

func fromPlanInternalCandidateEntity(
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
