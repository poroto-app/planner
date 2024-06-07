package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"time"
)

func (s Service) ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeIdToReplace string) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}

	planToUpdate := planCandidate.GetPlan(planId)
	if planToUpdate == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	// 入れ替え対象となる場所を取得
	placeToBeReplaced := planToUpdate.GetPlace(placeIdToBeReplaced)
	if placeToBeReplaced == nil {
		return nil, fmt.Errorf("place to be replaced not found: %v\n", placeIdToBeReplaced)
	}

	// 指定された場所がすでにプランに含まれている場合は何もしない
	if planToUpdate.GetPlace(placeIdToReplace) != nil {
		return nil, fmt.Errorf("place to replace already exists: %v\n", placeIdToReplace)
	}

	placeToReplace, err := s.placeRepository.Find(ctx, placeIdToReplace)
	if err != nil {
		return nil, fmt.Errorf("error while fetching place: %v\n", err)
	}

	if placeToReplace == nil {
		return nil, fmt.Errorf("place to replace not found: %v\n", placeIdToReplace)
	}

	if err := s.planCandidateRepository.ReplacePlace(ctx, planCandidateId, planId, placeIdToBeReplaced, *placeToReplace); err != nil {
		return nil, fmt.Errorf("error while replacing place: %v\n", err)
	}

	planCandidateUpdated, err := s.planCandidateRepository.Find(ctx, planCandidateId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}

	planUpdated := planCandidateUpdated.GetPlan(planId)
	if planUpdated == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	return planUpdated, nil
}
