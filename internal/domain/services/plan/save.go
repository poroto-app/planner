package plan

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/utils"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) SavePlanFromPlanCandidate(ctx context.Context, planCandidateId string, planId string, authToken *string) (*models.Plan, error) {
	// プラン候補から対応するプランを取得
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId, time.Now())
	if err != nil {
		return nil, err
	}

	planToSave, ok := array.Find(planCandidate.Plans, func(plan models.Plan) bool {
		return plan.Id == planId
	})
	if !ok {
		return nil, fmt.Errorf("plan(%v) not found in plan candidate(%v)", planId, planCandidateId)
	}

	// 冪等性を保つために、既存のプランを取得してから保存する
	planSaved, err := s.planRepository.Find(ctx, planId)
	if err != nil {
		// ログに出力するが、エラーは返さない
		s.logger.Warn(
			"error while finding plan",
			zap.String("planId", planId),
			zap.Error(err),
		)
	}

	if planSaved != nil {
		s.logger.Debug(
			"plan already exists. skip saving plan",
			zap.String("planId", planId),
		)
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
	if err := s.planRepository.Save(ctx, &planToSave); err != nil {
		return nil, err
	}

	return &planToSave, nil
}
