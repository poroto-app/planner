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
	TimeInMinutes   int      `firestore:"time_in_minutes"`
	PlaceIdsOrdered []string `firestore:"place_ids_ordered"`
}

func toPlanInCandidateEntity(
	id string,
	name string,
	places []models.Place,
	timeInMinutes uint,
) PlanInCandidateEntity {
	ps := make([]PlaceEntity, len(places))
	placeIdsOrdered := make([]string, len(places))

	for i, place := range places {
		ps[i] = ToPlaceEntity(place)
		placeIdsOrdered[i] = place.Id
	}

	return PlanInCandidateEntity{
		Id:              id,
		Name:            name,
		Places:          ps,
		TimeInMinutes:   int(timeInMinutes),
		PlaceIdsOrdered: placeIdsOrdered,
	}
}

func fromPlanInCandidateEntity(
	id string,
	name string,
	places []PlaceEntity,
	timeInMinutes int,
	placeIdsOrdered []string,
) models.Plan {
	return fromPlanEntity(
		id,
		name,
		places,
		timeInMinutes,
		placeIdsOrdered,
	)
}
