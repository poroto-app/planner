package plan

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/user"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	planRepository          repository.PlanRepository
	planCandidateRepository repository.PlanCandidateRepository
	userService             *user.Service
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

	userService, err := user.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing user service: %v", err)
	}

	return &Service{
		planRepository:          planRepository,
		planCandidateRepository: planCandidateRepository,
		userService:             userService,
	}, err
}
