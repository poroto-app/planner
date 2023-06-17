package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

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
	ps := make([]models.Place, len(places))
	for i, place := range places {
		ps[i] = FromPlaceEntity(place)
	}

	return models.Plan{
		Id:            id,
		Name:          name,
		Places:        ps,
		TimeInMinutes: uint(timeInMinutes),
	}
}
