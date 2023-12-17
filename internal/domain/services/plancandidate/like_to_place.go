package plancandidate

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) LikeToPlaceInPlanCandidate(
	ctx context.Context,
	planCandidateId string,
	placeId string,
	like bool,
) (*models.PlanCandidate, error) {

	err := s.planCandidateRepository.UpdateLikeToPlaceInPlanCandidate(ctx, planCandidateId, placeId, like)
	if err != nil {
		return nil, fmt.Errorf("error while updating like to place in plan candidate: %v", err)
	}
	placeCandidate, err := s.FindPlanCandidate(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate after updating: %v", err)
	}

	return placeCandidate, nil
}
