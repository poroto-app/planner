package repository

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanRepository interface {
	Save(ctx context.Context, plan *models.Plan) error
	Find(ctx context.Context, planId string) (*models.Plan, error)
}
