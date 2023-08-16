package plan

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	planRepository          repository.PlanRepository
	planCandidateRepository repository.PlanCandidateRepository
}

func NewService(ctx context.Context) (*Service, error) {
	planRepository, err := firestore.NewPlanRepository(ctx)
	if err != nil {
		return nil, err
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		planRepository:          planRepository,
		planCandidateRepository: planCandidateRepository,
	}, err
}
