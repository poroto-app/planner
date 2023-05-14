package services

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s PlanService) CachePlanCandidate(session string, plans []models.Plan) error {
	return s.planCandidateRepository.Save(&models.PlanCandidate{
		Id:        session,
		Plans:     plans,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
}

func (s PlanService) FindPlanCandidate(planCandidateId string) (*models.PlanCandidate, error) {
	return s.planCandidateRepository.Find(planCandidateId)
}
