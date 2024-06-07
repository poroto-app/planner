package plancandidate

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// DeleteExpiredPlanCandidates は期限切れのプラン候補を削除する
// それに伴いプラン候補に紐づくデータ（検索結果のキャッシュ等）も削除する
func (s Service) DeleteExpiredPlanCandidates(ctx context.Context, expiresAt time.Time) error {
	s.logger.Info(
		"Fetching expired plan candidates",
		zap.Time("expiresAt", expiresAt),
	)
	expiredPlanCandidateSetIds, err := s.planCandidateRepository.FindExpiredBefore(ctx, expiresAt)
	if err != nil {
		return fmt.Errorf("error while finding expired plan candidates: %v", err)
	}

	if len(*expiredPlanCandidateSetIds) == 0 {
		s.logger.Info("No expired plan candidates found")
		return nil
	}

	s.logger.Info(
		"Found expired plan candidate sets",
		zap.Int("count", len(*expiredPlanCandidateSetIds)),
	)

	// プラン候補を削除
	s.logger.Info(
		"Deleting expired plan candidates",
		zap.Int("count", len(*expiredPlanCandidateSetIds)),
	)
	if err := s.planCandidateRepository.DeleteAll(ctx, *expiredPlanCandidateSetIds); err != nil {
		return fmt.Errorf("error while deleting expired plan candidates: %v", err)
	}
	s.logger.Info(
		"Successfully deleted expired plan candidates",
		zap.Int("count", len(*expiredPlanCandidateSetIds)),
	)

	return nil
}
