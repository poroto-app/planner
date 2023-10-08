package plangen

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s Service) CreatePlanFromPlace(
	ctx context.Context,
	createPlanSessionId string,
	placeId string,
) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate")
	}

	// TODO: ユーザーの興味等を保存しておいて、それを反映させる
	placesSearched, err := s.placeSearchResultRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	var placeStart *places.Place
	for _, place := range placesSearched {
		if place.PlaceID == placeId {
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
			locationStart:                placeStart.Location.ToGeoLocation(),
			placeStart:                   *placeStart,
			places:                       placesSearched,
			freeTime:                     nil, // TODO: freeTimeの項目を保存し、それを反映させる
			createBasedOnCurrentLocation: planCandidate.MetaData.CreatedBasedOnCurrentLocation,
			shouldOpenWhileTraveling:     false, // 場所を検索してプランを作成した場合、必ずしも今すぐ行くとは限らない
		},
	)
	if err != nil {
		return nil, err
	}

	plansCreated := s.createPlanData(ctx, createPlanSessionId, CreatePlanParams{
		locationStart: placeStart.Location.ToGeoLocation(),
		placeStart:    *placeStart,
		places:        planPlaces,
	})
	if len(plansCreated) == 0 {
		return nil, fmt.Errorf("no plan created")
	}

	plan := plansCreated[0]

	if _, err = s.planCandidateRepository.AddPlan(ctx, createPlanSessionId, &plan); err != nil {
		return nil, err
	}

	return &plan, nil
}
