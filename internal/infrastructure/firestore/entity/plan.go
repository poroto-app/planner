package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanEntity は保存されたプランを示す
// GeoHash はプランの最初の場所のGeoHashを示す（プランは小さい範囲で作られるため、どこをとってもあまり変わらない）
// TimeInMinutes MEMO: Firestoreではuintをサポートしていないため，intにしている
type PlanEntity struct {
	Id              string               `firestore:"id"`
	Name            string               `firestore:"name"`
	Places          []PlaceEntity        `firestore:"places"`
	GeoHash         *string              `firestore:"geohash,omitempty"`
	PlaceIdsOrdered []string             `firestore:"place_ids_ordered"`
	TimeInMinutes   int                  `firestore:"time_in_minutes"`
	Transitions     *[]TransitionsEntity `firestore:"transitions,omitempty"`
	CreatedAt       time.Time            `firestore:"created_at,omitempty,serverTimestamp"`
	UpdatedAt       time.Time            `firestore:"updated_at,omitempty"`
}

func ToPlanEntity(plan models.Plan) PlanEntity {
	places := make([]PlaceEntity, len(plan.Places))
	placeIdsOrdered := make([]string, len(places))

	for i, place := range plan.Places {
		places[i] = ToPlaceEntity(place)
		placeIdsOrdered[i] = place.Id
	}

	var geohash *string
	if len(plan.Places) > 0 {
		value := plan.Places[0].Location.GeoHash()
		geohash = &value
	}

	return PlanEntity{
		Id:              plan.Id,
		Name:            plan.Name,
		Places:          places,
		GeoHash:         geohash,
		TimeInMinutes:   int(plan.TimeInMinutes),
		Transitions:     ToTransitionsEntities(plan.Transitions),
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
		entity.Transitions,
	)
}

func fromPlanEntity(
	id string,
	name string,
	places []PlaceEntity,
	timeInMinutes int,
	placeIdsOrdered []string,
	transitions *[]TransitionsEntity,
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
		Transitions:   FromTransitionEntities(transitions),
	}
}
