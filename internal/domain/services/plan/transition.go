package plan

import "poroto.app/poroto/planner/internal/domain/models"

func (s PlanService) addTransition(
	placesInPlan []models.Place,
	transitions []models.Transition,
	duration uint,
	createBasedOnCurrentLocation bool,
) []models.Transition {
	if len(placesInPlan) == 0 {
		return transitions
	}

	// 場所を指定して作成した場合は、最初の場所をFromにする
	if !createBasedOnCurrentLocation && len(placesInPlan) == 1 {
		return transitions
	}

	var fromPlaceId *string
	if createBasedOnCurrentLocation && len(placesInPlan) == 1 {
		// 現在地から作成した場合は最初のFromがnilになる
		fromPlaceId = nil
	} else {
		// そうでない場合は基準となる場所がFromになる
		fromPlaceId = &placesInPlan[len(placesInPlan)-2].Id
	}

	transitions = append(transitions, models.Transition{
		FromPlaceId: fromPlaceId,
		ToPlaceId:   placesInPlan[len(placesInPlan)-1].Id,
		Duration:    duration,
	})

	return transitions
}
