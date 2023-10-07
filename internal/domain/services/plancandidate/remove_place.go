package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

// RemovePlaceFromPlan プラン候補から場所を削除する
// planId に対応するプランが存在しない場合はエラーを返す
// 指定された場所をプランから除外すると、プランに含まれる場所が0になる場合はエラーを返す
func (s Service) RemovePlaceFromPlan(ctx context.Context, planCandidateId string, planId string, placeId string) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving plan candidate: %v", err)
	}

	plan := planCandidate.GetPlan(planId)
	if err != nil {
		return nil, fmt.Errorf("plan not found in plan candidate: %v", err)
	}

	// 少なくとも1つの場所がプランに含まれるようにする
	if len(plan.Places) <= 1 {
		return nil, fmt.Errorf("cannot remove last place from plan")
	}

	// プラン候補から場所を削除
	err = s.planCandidateRepository.RemovePlaceFromPlan(ctx, planCandidateId, planId, placeId)
	if err != nil {
		return nil, fmt.Errorf("error while removing place from plan candidate: %v", err)
	}

	// 更新後のプラン候補を取得
	planCandidate, err = s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving plan candidate: %v", err)
	}

	plan = planCandidate.GetPlan(planId)
	if err != nil {
		return nil, fmt.Errorf("plan not found in plan candidate: %v", err)
	}

	return plan, nil
}
