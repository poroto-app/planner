package repository

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateRepository interface {
	Save(cxt context.Context, planCandidate *models.PlanCandidate) error
	Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error)
	FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]models.PlanCandidate, error)
	AddPlan(ctx context.Context, planCandidateId string, plan *models.Plan) (*models.PlanCandidate, error)
	UpdatePlacesOrder(ctx context.Context, planId string, planCandidate string, placeIdsOrdered []string) (*models.Plan, error)
	DeleteAll(ctx context.Context, planCandidateIds []string) error
}
