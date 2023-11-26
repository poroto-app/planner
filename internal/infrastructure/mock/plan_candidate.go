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

func (p PlanRepository) Create(ctx context.Context, planCandidateId string, expiresAt time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]string, error) {
	var values []string
	for _, value := range p.Data {
		if !value.ExpiresAt.After(expiresAt) {
			values = append(values, value.Id)
		}
	}

	return &values, nil
}

func (p PlanRepository) AddPlan(ctx context.Context, planCandidateId string, plan ...models.Plan) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) AddPlaceToPlan(ctx context.Context, planCandidateId string, planId string, previousPlaceId string, place models.Place) error {
	// TODO: implement me
	panic("implement me")
}

func (p PlanRepository) RemovePlaceFromPlan(ctx context.Context, planCandidateId string, planId string, placeId string) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) UpdatePlacesOrder(ctx context.Context, planId string, planCandidate string, placeIdsOrdered []string) (*models.Plan, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) UpdatePlanCandidateMetaData(ctx context.Context, planCandidateId string, meta models.PlanCandidateMetaData) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeToReplace models.Place) error {
	// TODO implement me
	panic("implement me")
}

func (p PlanRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	for _, id := range planCandidateIds {
		delete(p.Data, id)
	}
	return nil
}
