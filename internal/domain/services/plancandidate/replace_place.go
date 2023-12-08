package plancandidate

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeIdToReplace string) (*models.Plan, error) {
	s.logger.Info(
		"start replacing place",
		zap.String("planCandidateId", planCandidateId),
		zap.String("planId", planId),
		zap.String("placeIdToBeReplaced", placeIdToBeReplaced),
		zap.String("placeIdToReplace", placeIdToReplace),
	)
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}
	s.logger.Info(
		"succeeded fetching plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)

	planToUpdate := planCandidate.GetPlan(planId)
	if planToUpdate == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	s.logger.Info(
		"start fetching searched places",
		zap.String("planCandidateId", planCandidateId),
	)
	places, err := s.placeService.FetchSearchedPlaces(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}
	s.logger.Info(
		"succeeded fetching searched places",
		zap.String("planCandidateId", planCandidateId),
	)

	s.logger.Info(
		"start fetching place to be replaced",
		zap.String("planCandidateId", planCandidateId),
		zap.String("placeIdToBeReplaced", placeIdToBeReplaced),
	)
	placeToBeReplaced := planToUpdate.GetPlace(placeIdToBeReplaced)
	if placeToBeReplaced == nil {
		return nil, fmt.Errorf("place to be replaced not found: %v\n", placeIdToBeReplaced)
	}
	s.logger.Info(
		"succeeded fetching place to be replaced",
		zap.String("planCandidateId", planCandidateId),
		zap.String("placeIdToBeReplaced", placeIdToBeReplaced),
	)

	// 指定された場所がすでにプランに含まれている場合は何もしない
	if planToUpdate.GetPlace(placeIdToReplace) != nil {
		return nil, fmt.Errorf("place to replace already exists: %v\n", placeIdToReplace)
	}

	var placeToReplace *models.Place
	for _, place := range places {
		if place.Id == placeIdToReplace {
			placeToReplace = &place
			break
		}
	}
	if placeToReplace == nil {
		return nil, fmt.Errorf("place to replace not found: %v\n", placeIdToReplace)
	}

	s.logger.Info(
		"start fetching photos and reviews for places",
		zap.String("planCandidateId", planCandidateId),
	)
	if err := s.planCandidateRepository.ReplacePlace(ctx, planCandidateId, planId, placeIdToBeReplaced, *placeToReplace); err != nil {
		return nil, fmt.Errorf("error while replacing place: %v\n", err)
	}
	s.logger.Info(
		"succeeded fetching photos and reviews for places",
		zap.String("planCandidateId", planCandidateId),
	)

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
