package repository

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceInPlanCandidateRepository interface {
	Save(ctx context.Context, planCandidateId string, place models.PlaceInPlanCandidate) error

	SavePlaces(ctx context.Context, planCandidateId string, places []models.PlaceInPlanCandidate) error

	FindByPlanCandidateId(ctx context.Context, planCandidateId string) (*[]models.PlaceInPlanCandidate, error)

	DeleteByPlanCandidateId(ctx context.Context, planCandidateId string) error
}
