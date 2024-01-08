package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) SaveSearchedPlaces(ctx context.Context, planCandidateId string, googlePlaces []models.GooglePlace) ([]models.Place, error) {
	// models.Google を保存し，models.Place を取得する
	places, err := s.placeRepository.SavePlacesFromGooglePlaces(ctx, googlePlaces...)
	if err != nil {
		return nil, fmt.Errorf("error while saving places from google place: %v\n", err)
	}

	// PlanCandidate と検索された場所の紐付けを行う
	var placeIds []string
	for _, place := range *places {
		placeIds = append(placeIds, place.Id)
	}
	if err := s.planCandidateRepository.AddSearchedPlacesForPlanCandidate(ctx, planCandidateId, placeIds); err != nil {
		return nil, fmt.Errorf("error while adding searched places for plan candidate: %v\n", err)
	}

	return *places, nil
}

func (s Service) FetchSearchedPlaces(ctx context.Context, planCandidateId string) ([]models.Place, error) {
	places, err := s.placeRepository.FindByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching searched places for plan candidate: %v\n", err)
	}

	return places, nil
}
