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
	meta models.PlanCandidateMetaData,
) error {
	return s.planCandidateRepository.Save(ctx, &models.PlanCandidate{
		Id:        session,
		Plans:     plans,
		MetaData:  meta,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
}
