package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanEntity struct {
	Id     string        `firestore:"id"`
	Name   string        `firestore:"name"`
	Places []PlaceEntity `firestore:"places"`
	// Updateの範囲を減らすため，Placesの順番をID配列の順で管理している
	PlaceIdsOrdered []string `firestore:"place_ids_ordered"`
	// MEMO: Firestoreではuintをサポートしていないため，intにしている
	TimeInMinutes int       `firestore:"time_in_minutes"`
	CreatedAt     time.Time `firestore:"created_at,omitempty,serverTimestamp"`
	UpdatedAt     time.Time `firestore:"updated_at,omitempty"`
}

func ToPlanEntity(plan models.Plan) PlanEntity {
	places := make([]PlaceEntity, len(plan.Places))
	placeIdsOrdered := make([]string, len(places))

	for i, place := range plan.Places {
		places[i] = ToPlaceEntity(place)
		placeIdsOrdered[i] = place.Id
	}

	return PlanEntity{
		Id:              plan.Id,
		Name:            plan.Name,
		Places:          places,
		TimeInMinutes:   int(plan.TimeInMinutes),
		PlaceIdsOrdered: placeIdsOrdered,
	}
}

func FromPlanEntity(entity PlanEntity) models.Plan {
	return fromPlanEntity(
		entity.Id,
		entity.Name,
		entity.Places,
		entity.TimeInMinutes,
		entity.PlaceIdsOrdered,
	)
}

func fromPlanEntity(
	id string,
	name string,
	places []PlaceEntity,
	timeInMinutes int,
	placeIdsOrdered []string,
) models.Plan {
	// placeIdsOrdered：プレイスの順序を指定するプレイスのID配列
	// データベースモデルからドメインモデルに変換する際にプレイスの順序を並び替える
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
	}
}
