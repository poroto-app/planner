package plangen

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

// CreatePlanFromPlace 指定した場所を起点としてプランを作成する
func (s Service) CreatePlanFromPlace(
	ctx context.Context,
	createPlanSessionId string,
	placeId string,
) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, createPlanSessionId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate")
	}

	placeStart, err := s.placeRepository.Find(ctx, placeId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching place")
	}

	if placeStart == nil {
		return nil, fmt.Errorf("place not found")
	}

	placesNearby, err := s.placeSearchService.SearchNearbyPlaces(ctx, placesearch.SearchNearbyPlacesInput{
		Location:           placeStart.Location,
		PlanCandidateSetId: &createPlanSessionId,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching nearby places")
	}

	var categoryNamesRejected []string
	if planCandidate.MetaData.CategoriesRejected != nil {
		for _, category := range *planCandidate.MetaData.CategoriesRejected {
			categoryNamesRejected = append(categoryNamesRejected, category.Name)
		}
	}

	// TODO: ユーザーの興味等を保存しておいて、それを反映させる
	planPlaces, err := s.CreatePlanPlaces(CreatePlanPlacesInput{
		PlanCandidateId:       createPlanSessionId,
		LocationStart:         placeStart.Location,
		PlaceStart:            *placeStart,
		Places:                placesNearby,
		CategoryNamesDisliked: &categoryNamesRejected,
		FreeTime:              planCandidate.MetaData.FreeTime,
	})
	if err != nil {
		return nil, err
	}

	plansCreated := s.createPlanData(ctx, createPlanSessionId, CreatePlanParams{
		LocationStart: placeStart.Location,
		PlaceStart:    *placeStart,
		Places:        planPlaces,
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
