package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/services/user"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type FindPlanCandidateSetInput struct {
	PlanCandidateSetId string
	UserId             *string
	FirebaseAuthToken  *string
}

func (s Service) Find(ctx context.Context, input FindPlanCandidateSetInput) (*models.PlanCandidateSet, error) {
	planCandidateSet, err := s.planCandidateRepository.Find(ctx, input.PlanCandidateSetId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error finding plan candidate: %w", err)
	}

	// ログインしている場合は、そのユーザーがいいねした場所を取得する
	if input.UserId != nil && input.FirebaseAuthToken != nil {
		checkAuthResult, err := s.userService.CheckUserAuthState(ctx, user.CheckUserAuthStateInput{
			UserId:            *input.UserId,
			FirebaseAuthToken: *input.FirebaseAuthToken,
		})
		if err != nil {
			return nil, err
		}

		if !checkAuthResult.IsAuthenticated {
			return nil, fmt.Errorf("user is not authorized")
		}

		likePlaces, err := s.placeRepository.FindLikePlacesByUserId(ctx, *input.UserId)
		if err != nil {
			return nil, fmt.Errorf("error finding like places: %w", err)
		}

		planCandidateSet.LikedPlaceIds = array.Map(*likePlaces, func(place models.Place) string {
			return place.Id
		})
	}

	return planCandidateSet, nil
}
