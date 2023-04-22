package repository

import "poroto.app/poroto/planner/internal/domain/models"

type PlanRepository interface {
	Save(plan *models.Plan) error
	Find(planId *models.Plan) (*models.Plan, error)
}
