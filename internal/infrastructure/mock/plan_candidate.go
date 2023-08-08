package mock

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanRepository struct {
	Data map[string]models.PlanCandidate
}

func NewPlanCandidateRepository(data map[string]models.PlanCandidate) *PlanRepository {
	return &PlanRepository{
		Data: data,
	}
}

func (p PlanRepository) Save(cxt context.Context, planCandidate *models.PlanCandidate) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]models.PlanCandidate, error) {
	var values []models.PlanCandidate
	for _, value := range p.Data {
		if value.ExpiresAt.Before(expiresAt) {
			values = append(values, value)
		}
	}

	return &values, nil
}

func (p PlanRepository) AddPlan(ctx context.Context, planCandidateId string, plan *models.Plan) (*models.PlanCandidate, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) UpdatePlacesOrder(ctx context.Context, planId string, planCandidate string, placeIdsOrdered []string) (*models.Plan, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	for _, id := range planCandidateIds {
		delete(p.Data, id)
	}
	return nil
}
