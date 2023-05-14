package inmemory

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

var data = make(map[string]*models.PlanCandidate)

type PlanCandidateInMemoryRepository struct {
}

func NewPlanCandidateRepository() *PlanCandidateInMemoryRepository {
	return &PlanCandidateInMemoryRepository{}
}

func (p *PlanCandidateInMemoryRepository) Save(planCandidate *models.PlanCandidate) error {
	data[planCandidate.Id] = planCandidate
	return nil
}

func (p *PlanCandidateInMemoryRepository) Find(planCandidateId *models.PlanCandidate) (*models.PlanCandidate, error) {
	if candidate, ok := data[planCandidateId.Id]; ok {
		return candidate, nil
	}
	return nil, nil
}
