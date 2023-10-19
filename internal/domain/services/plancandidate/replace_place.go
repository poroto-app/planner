package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeIdToReplace string) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}

	planToUpdate := planCandidate.GetPlan(planId)
	if planToUpdate == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	placesSearched, err := s.placeSearchResultRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}

	placeToBeReplaced := planToUpdate.GetPlace(placeIdToBeReplaced)
	if placeToBeReplaced == nil {
		return nil, fmt.Errorf("place to be replaced not found: %v\n", placeIdToBeReplaced)
	}

	var placeToReplace *models.Place
	for _, place := range placesSearched {
		// TODO: PlaceRepositoryを用いて、Planner APIが指定したPlaceIdで取得できるようにする
		if place.PlaceId == placeIdToReplace {
			*placeToReplace = place.ToPlace()
			break
		}
	}
	if placeToReplace == nil {
		return nil, fmt.Errorf("place to replace not found: %v\n", placeIdToReplace)
	}

	if err := s.planCandidateRepository.ReplacePlace(ctx, planCandidateId, planId, placeIdToBeReplaced, *placeToReplace); err != nil {
		return nil, fmt.Errorf("error while replacing place: %v\n", err)
	}

	planCandidateUpdated, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}

	planUpdated := planCandidateUpdated.GetPlan(planId)
	if planUpdated == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	return planUpdated, nil
}
