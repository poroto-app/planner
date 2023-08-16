package plan

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/plangen"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type PlanService struct {
	placesApi                   places.PlacesApi
	planRepository              repository.PlanRepository
	planCandidateRepository     repository.PlanCandidateRepository
	placeSearchResultRepository repository.PlaceSearchResultRepository
	planGeneratorService        plangen.Service
}

func NewPlanService(ctx context.Context) (*PlanService, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	planRepository, err := firestore.NewPlanRepository(ctx)
	if err != nil {
		return nil, err
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	placeSearchResultRepository, err := firestore.NewPlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, err
	}

	planGeneratorService, err := plangen.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan generator service: %v", err)
	}

	return &PlanService{
		placesApi:                   *placesApi,
		planRepository:              planRepository,
		planCandidateRepository:     planCandidateRepository,
		placeSearchResultRepository: placeSearchResultRepository,
		planGeneratorService:        *planGeneratorService,
	}, err
}

func (s PlanService) CreatePlanFromPlace(
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

	planCreated, err := s.planGeneratorService.CreatePlan(
		ctx,
		placeStart.Location.ToGeoLocation(),
		*placeStart,
		placesSearched,
		// TODO: freeTimeの項目を保存し、それを反映させる
		nil,
		planCandidate.CreatedBasedOnCurrentLocation,
	)
	if err != nil {
		return nil, err
	}

	if _, err = s.planCandidateRepository.AddPlan(ctx, createPlanSessionId, planCreated); err != nil {
		return nil, err
	}

	return planCreated, nil
}
