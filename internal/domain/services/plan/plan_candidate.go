package plan

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s PlanService) CachePlanCandidate(ctx context.Context, session string, plans []models.Plan, createdBasedOnCurrentLocation bool) error {
	return s.planCandidateRepository.Save(ctx, &models.PlanCandidate{
		Id:                            session,
		Plans:                         plans,
		CreatedBasedOnCurrentLocation: createdBasedOnCurrentLocation,
		ExpiresAt:                     time.Now().Add(7 * 24 * time.Hour),
	})
}

func (s PlanService) FindPlanCandidate(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	return s.planCandidateRepository.Find(ctx, planCandidateId)
}

func (s PlanService) UpdatePlacesOrderPlanCandidate(ctx context.Context, planId string, planCandidate *models.PlanCandidate, placeIdsOrdered []string) (*models.Plan, error) {
	return s.planCandidateRepository.UpdatePlacesOrder(ctx, planId, planCandidate, placeIdsOrdered)
}
