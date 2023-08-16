package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanInCandidateEntity PlanCandidateEntityに含まれるPlan
// MEMO: PlanEntityを用いると、CreatedAtとUpdatedAtが含まれてしまうため、別の構造体を利用している
type PlanInCandidateEntity struct {
	Id          string               `firestore:"id"`
	Name        string               `firestore:"name"`
	Places      []PlaceEntity        `firestore:"places"`
	Transitions *[]TransitionsEntity `firestore:"transitions,omitempty"`
	// MEMO: Firestoreではuintをサポートしていないため，intにしている
	TimeInMinutes   int      `firestore:"time_in_minutes"`
	PlaceIdsOrdered []string `firestore:"place_ids_ordered"`
}

func ToPlanInCandidateEntity(plan models.Plan) PlanInCandidateEntity {
	ps := make([]PlaceEntity, len(plan.Places))
	placeIdsOrdered := make([]string, len(plan.Places))

	for i, place := range plan.Places {
		ps[i] = ToPlaceEntity(place)
		placeIdsOrdered[i] = place.Id
	}

	return PlanInCandidateEntity{
		Id:              plan.Id,
		Name:            plan.Name,
		Places:          ps,
		PlaceIdsOrdered: placeIdsOrdered,
		TimeInMinutes:   int(plan.TimeInMinutes),
		Transitions:     ToTransitionsEntities(plan.Transitions),
	}
}

func fromPlanInCandidateEntity(
	id string,
	name string,
	places []PlaceEntity,
	timeInMinutes int,
	placeIdsOrdered []string,
	transitions *[]TransitionsEntity,
) models.Plan {
	return fromPlanEntity(
		id,
		name,
		places,
		timeInMinutes,
		placeIdsOrdered,
		transitions,
	)
}
