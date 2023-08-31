package plan

import (
	"context"
	"fmt"
	"log"
	"poroto.app/poroto/planner/internal/domain/utils"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) SavePlanFromPlanCandidate(ctx context.Context, planCandidateId string, planId string, authToken *string) (*models.Plan, error) {
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
	planSaved, err := s.planRepository.Find(ctx, planId)
	if err != nil {
		// ログに出力するが、エラーは返さない
		log.Println(fmt.Errorf("error while finding plan(%v): %v", planId, err))
	}

	if planSaved != nil {
		log.Printf("plan(%v) already exists. skip saving plan", planId)
		return planSaved, nil
	}

	// ユーザー情報を取得
	if authToken != nil {
		user, err := s.userService.FindByFirebaseIdToken(ctx, *authToken)
		if err != nil {
			return nil, fmt.Errorf("error while getting user from firebase id token: %v", err)
		}

		if user == nil {
			return nil, fmt.Errorf("user not found")
		}

		planToSave.AuthorId = utils.StrPointer(user.Id)
	}

	// プランを保存
	if err := s.planRepository.Save(ctx, planToSave); err != nil {
		return nil, err
	}

	return planToSave, nil
}
