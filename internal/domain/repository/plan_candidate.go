package repository

import "poroto.app/poroto/planner/internal/domain/models"

type PlanCandidateRepository interface {
	Save(planCandidate *models.PlanCandidate) error
	Find(planCandidateId string) (*models.PlanCandidate, error)
}
