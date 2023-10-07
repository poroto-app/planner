package entity

import (
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanInCandidateEntity PlanCandidateEntityに含まれるPlan
// PlaceIdsOrdered は Places の順番を管理する（Places配列を書き換えて更新すると、更新の量が多くなるため）
// MEMO: PlanEntityを用いると、CreatedAtとUpdatedAtが含まれてしまうため、別の構造体を利用している
type PlanInCandidateEntity struct {
	Id              string        `firestore:"id"`
	Name            string        `firestore:"name"`
	Places          []PlaceEntity `firestore:"places"`
	PlaceIdsOrdered []string      `firestore:"place_ids_ordered"`
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
	}
}

func fromPlanInCandidateEntity(
	id string,
	name string,
	places []PlaceEntity,
	placeIdsOrdered []string,
	timeInMinutes int,
) (*models.Plan, error) {
	placesOrdered := make([]models.Place, len(places))

	if !validatePlanInCandidateEntity(places, placeIdsOrdered) {
		return nil, fmt.Errorf("the property of placeIdsOrdered is invalid")
	}

	// 整合性がある場合，指定された順番でドメインモデルに変換
	for i, placeIdOrdered := range placeIdsOrdered {
		for _, place := range places {
			if place.Id == placeIdOrdered {
				placesOrdered[i] = FromPlaceEntity(place)
			}
		}
	}

	return &models.Plan{
		Id:            id,
		Name:          name,
		Places:        placesOrdered,
		TimeInMinutes: uint(timeInMinutes),
	}, nil
}

// validatePlanInCandidateEntity はプラン候補内プランの場所一覧と順序指定のID配列の整合性をチェックする
func validatePlanInCandidateEntity(places []PlaceEntity, placeIdsOrdered []string) bool {
	// 順序指定ID配列の数が正しいかどうか　を確認
	if len(places) != len(placeIdsOrdered) {
		return false
	}

	// 順序指定ID配列の中に重複がないか，実在するPlace.Idに一致するか　を確認
	placeIncluded := make(map[string]PlaceEntity)
	for _, placeIdOrdered := range placeIdsOrdered {
		for _, place := range places {
			if place.Id == placeIdOrdered {
				placeIncluded[place.Id] = place
			}
		}
	}

	// ID配列内に重複がない and 順序指定ID配列のIDが正当なものである 場合　Trueを返す
	return len(placeIncluded) == len(places)
}
