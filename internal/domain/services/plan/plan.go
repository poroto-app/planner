package plan

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type PlanService struct {
	planRepository          repository.PlanRepository
	planCandidateRepository repository.PlanCandidateRepository
}

func NewPlanService(ctx context.Context) (*PlanService, error) {
	planRepository, err := firestore.NewPlanRepository(ctx)
	if err != nil {
		return nil, err
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &PlanService{
		planRepository:          planRepository,
		planCandidateRepository: planCandidateRepository,
	}, err
}
