package plancandidate

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) SavePlanCandidate(
	ctx context.Context,
	session string,
	plans []models.Plan,
	createdBasedOnCurrentLocation bool,
) error {
	return s.planCandidateRepository.Save(ctx, &models.PlanCandidate{
		Id:                            session,
		Plans:                         plans,
		CreatedBasedOnCurrentLocation: createdBasedOnCurrentLocation,
		ExpiresAt:                     time.Now().Add(7 * 24 * time.Hour),
	})
}
