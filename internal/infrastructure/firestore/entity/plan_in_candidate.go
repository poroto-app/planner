package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanInCandidateEntity PlanCandidateEntityに含まれるPlan
// MEMO: PlanEntityを用いると、CreatedAtとUpdatedAtが含まれてしまうため、別の構造体を利用している
type PlanInCandidateEntity struct {
	Id              string               `firestore:"id"`
	Name            string               `firestore:"name"`
	Places          []PlaceEntity        `firestore:"places"`
	PlaceIdsOrdered []string             `firestore:"place_ids_ordered"`
	Transitions     *[]TransitionsEntity `firestore:"transitions,omitempty"`
	// MEMO: Firestoreではuintをサポートしていないため，intにしている
	TimeInMinutes int `firestore:"time_in_minutes"`
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
	placeIdsOrdered []string,
	timeInMinutes int,
	transitions *[]TransitionsEntity,
) models.Plan {
	placesOrdered := make([]models.Place, len(places))
	for i, placeIdOrdered := range placeIdsOrdered {
		for _, place := range places {
			if place.Id == placeIdOrdered {
				placesOrdered[i] = FromPlaceEntity(place)
			}
		}
	}

	return models.Plan{
		Id:            id,
		Name:          name,
		Places:        placesOrdered,
		TimeInMinutes: uint(timeInMinutes),
		Transitions:   FromTransitionEntities(transitions),
	}
}
