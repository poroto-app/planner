package plangen

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) CreatePlanFromPlace(
	ctx context.Context,
	createPlanSessionId string,
	placeId string,
) (*models.Plan, error) {
	// TODO: ユーザーの興味等を保存しておいて、それを反映させる
	places, err := s.placeService.FetchSearchedPlaces(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	var placeStart *models.Place
	for _, place := range places {
		if place.Id == placeId {
			placeStart = &place
			break
		}
	}

	if placeStart == nil {
		return nil, fmt.Errorf("place not found")
	}

	planPlaces, err := s.createPlanPlaces(
		ctx,
		CreatePlanPlacesParams{
			PlanCandidateId:          createPlanSessionId,
			LocationStart:            placeStart.Location,
			PlaceStart:               *placeStart,
			Places:                   places,
			FreeTime:                 nil,   // TODO: freeTimeの項目を保存し、それを反映させる
			ShouldOpenWhileTraveling: false, // 場所を検索してプランを作成した場合、必ずしも今すぐ行くとは限らない
		},
	)
	if err != nil {
		return nil, err
	}

	plansCreated := s.createPlanData(ctx, createPlanSessionId, CreatePlanParams{
		locationStart: placeStart.Location,
		placeStart:    *placeStart,
		places:        planPlaces,
	})
	if len(plansCreated) == 0 {
		return nil, fmt.Errorf("no plan created")
	}

	plan := plansCreated[0]

	if err = s.planCandidateRepository.AddPlan(ctx, createPlanSessionId, plan); err != nil {
		return nil, err
	}

	return &plan, nil
}
