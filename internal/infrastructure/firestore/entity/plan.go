package entity

import "poroto.app/poroto/planner/internal/domain/models"

type PlanEntity struct {
	Id            string        `firestore:"id"`
	Name          string        `firestore:"name"`
	Places        []PlaceEntity `firestore:"places"`
	TimeInMinutes uint          `firestore:"time_in_minutes"`
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
		TimeInMinutes: plan.TimeInMinutes,
	}
}

func FromPlanEntity(entity PlanEntity) models.Plan {
	places := make([]models.Place, len(entity.Places))
	for i, place := range entity.Places {
		places[i] = FromPlaceEntity(place)
	}

	return models.Plan{
		Id:            entity.Id,
		Name:          entity.Name,
		Places:        places,
		TimeInMinutes: entity.TimeInMinutes,
	}
}
