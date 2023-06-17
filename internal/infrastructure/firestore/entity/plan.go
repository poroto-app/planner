package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanEntity struct {
	Id     string        `firestore:"id"`
	Name   string        `firestore:"name"`
	Places []PlaceEntity `firestore:"places"`
	// MEMO: Firestoreではuintをサポートしていないため，intにしている
	TimeInMinutes int       `firestore:"time_in_minutes"`
	CreatedAt     time.Time `firestore:"created_at,omitempty,serverTimestamp"`
	UpdatedAt     time.Time `firestore:"updated_at,omitempty"`
}

func ToPlanEntity(plan models.Plan) PlanEntity {
	places := make([]PlaceEntity, len(plan.Places))
	for i, place := range plan.Places {
		places[i] = ToPlaceEntity(place)
	}

	return PlanEntity{
		Id:            plan.Id,
		Name:          plan.Name,
		Places:        places,
		TimeInMinutes: int(plan.TimeInMinutes),
	}
}

func FromPlanEntity(entity PlanEntity) models.Plan {
	return fromPlanEntity(
		entity.Id,
		entity.Name,
		entity.Places,
		entity.TimeInMinutes,
	)
}

func fromPlanEntity(
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
