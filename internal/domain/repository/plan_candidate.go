package repository

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateRepository interface {
	Save(cxt context.Context, planCandidate *models.PlanCandidate) error
	Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error)
	UpdatePlacesOrder(ctx context.Context, planId string, planCandidate *models.PlanCandidate, placeIdsOrdered []string) (*models.Plan, error)
}
