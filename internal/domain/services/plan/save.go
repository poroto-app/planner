package plan

import (
	"context"
	"fmt"
	"log"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s PlanService) SavePlanFromPlanCandidate(ctx context.Context, planCandidateId string, planId string) (*models.Plan, error) {
	// プラン候補から対応するプランを取得
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}

	var planToSave *models.Plan
	for _, plan := range planCandidate.Plans {
		if plan.Id == planId {
			planToSave = &plan
			break
		}
	}
	if planToSave == nil {
		return nil, fmt.Errorf("plan(%v) not found in plan candidate(%v)", planId, planCandidateId)
	}

	// 冪等性を保つために、既存のプランを取得してから保存する
	planSaved, err := s.planRepository.Find(planId)
	if err != nil {
		// ログに出力するが、エラーは返さない
		log.Println(fmt.Errorf("error while finding plan(%v): %v", planId, err))
	}

	if planSaved != nil {
		return planSaved, nil
	}

	// プランを保存
	if err := s.planRepository.Save(planToSave); err != nil {
		return nil, err
	}

	return planToSave, nil
}
