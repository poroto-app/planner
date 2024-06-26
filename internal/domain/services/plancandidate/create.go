package plancandidate

import (
	"context"
	"fmt"
	"time"
)

func (s Service) CreatePlanCandidateSet(
	ctx context.Context,
	planCandidateSetId string,
) error {
	if err := s.planCandidateRepository.Create(ctx, planCandidateSetId, time.Now().Add(7*24*time.Hour)); err != nil {
		return fmt.Errorf("error while creating plan candidate: %v\n", err)
	}
	return nil
}
