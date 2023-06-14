package inmemory

import "poroto.app/poroto/planner/internal/domain/models"

type PlanRepository struct {
	data map[string]*models.Plan
}

// TODO: Delete
// どこから初期化しても同じインスタンスが使えるようにする
var repository *PlanRepository

func NewPlanRepository() (*PlanRepository, error) {
	if repository == nil {
		repository = &PlanRepository{
			data: make(map[string]*models.Plan),
		}
	}
	return repository, nil
}

func (p *PlanRepository) Save(plan *models.Plan) error {
	p.data[plan.Id] = plan
	return nil
}

func (p *PlanRepository) Find(planId string) (*models.Plan, error) {
	for id, plan := range p.data {
		if id == planId {
			return plan, nil
		}
	}
	return nil, nil
}
