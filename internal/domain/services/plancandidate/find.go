package plancandidate

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) FindPlanCandidate(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	// TODO: ユーザーとして Like した場所を取得できるようにする
	return s.planCandidateRepository.Find(ctx, planCandidateId, time.Now())
}
