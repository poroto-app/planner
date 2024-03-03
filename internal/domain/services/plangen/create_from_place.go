package plangen

import (
	"context"
	"fmt"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) CreatePlanFromPlace(
	ctx context.Context,
	createPlanSessionId string,
	placeId string,
) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, createPlanSessionId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate")
	}

	// TODO: ユーザーの興味等を保存しておいて、それを反映させる
	places, err := s.placeSearchService.FetchSearchedPlaces(ctx, createPlanSessionId)
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

	var categoryNamesRejected []string
	if planCandidate.MetaData.CategoriesRejected != nil {
		for _, category := range *planCandidate.MetaData.CategoriesRejected {
			categoryNamesRejected = append(categoryNamesRejected, category.Name)
		}
	}

	planPlaces, err := s.CreatePlanPlaces(CreatePlanPlacesInput{
		PlanCandidateId:       createPlanSessionId,
		LocationStart:         placeStart.Location,
		PlaceStart:            *placeStart,
		Places:                places,
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
