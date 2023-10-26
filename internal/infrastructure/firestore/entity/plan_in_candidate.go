package entity

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanInCandidateEntity PlanCandidateEntityに含まれるPlan
// PlaceIdsOrdered は Places の順番を管理する（Places配列を書き換えて更新すると、更新の量が多くなるため）
// MEMO: PlanEntityを用いると、CreatedAtとUpdatedAtが含まれてしまうため、別の構造体を利用している
type PlanInCandidateEntity struct {
	Id              string   `firestore:"id"`
	Name            string   `firestore:"name"`
	PlaceIdsOrdered []string `firestore:"place_ids_ordered"`
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
		PlaceIdsOrdered: placeIdsOrdered,
		TimeInMinutes:   int(plan.TimeInMinutes),
	}
}

func FromPlanInCandidateEntity(
	id string,
	name string,
	places []models.PlaceInPlanCandidate,
	placeIdsOrdered []string,
	timeInMinutes int,
) (*models.Plan, error) {
	placesOrdered := make([]models.Place, len(placeIdsOrdered))
	for i, placeIdOrdered := range placeIdsOrdered {
		for _, place := range places {
			if place.Id == placeIdOrdered {
				placesOrdered[i] = place.ToPlace()
			}
		}
	}

	// 整合性の確認
	if err := validatePlaceInPlanCandidateEntity(placesOrdered, placeIdsOrdered); err != nil {
		return nil, fmt.Errorf("places in plan candidate is invalid: %w", err)
	}

	return &models.Plan{
		Id:            id,
		Name:          name,
		Places:        placesOrdered,
		TimeInMinutes: uint(timeInMinutes),
	}, nil
}

// validatePlaceInPlanCandidateEntity はプラン候補内プランの場所一覧と順序指定のID配列の整合性をチェックする
func validatePlaceInPlanCandidateEntity(places []models.Place, placeIdsOrdered []string) error {
	// 順序指定ID配列の数が正しいかどうか　を確認
	if len(places) != len(placeIdsOrdered) {
		return fmt.Errorf("the length of placeIdsOrdered is invalid")
	}

	placeIncluded := make(map[string]models.Place)
	for _, placeId := range placeIdsOrdered {
		for _, place := range places {
			if place.Id != placeId {
				continue
			}

			// 重複があるかどうかを確認
			if _, ok := placeIncluded[place.Id]; ok {
				return fmt.Errorf("place(%s) is duplicated", placeId)
			}

			placeIncluded[place.Id] = place
		}

		// 存在しない場所IDが指定されていないかどうかを確認
		if _, ok := placeIncluded[placeId]; !ok {
			return fmt.Errorf("place(%s) was not found", placeId)
		}
	}

	return nil
}
