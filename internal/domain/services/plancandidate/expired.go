package plancandidate

import (
	"context"
	"fmt"
	"log"
	"time"
)

// DeleteExpiredPlanCandidates は期限切れのプラン候補を削除する
// それに伴いプラン候補に紐づくデータ（検索結果のキャッシュ等）も削除する
func (s Service) DeleteExpiredPlanCandidates(ctx context.Context, expiresAt time.Time) error {
	log.Println("Fetching expired plan candidates")
	expiredPlanCandidateIds, err := s.planCandidateRepository.FindExpiredBefore(ctx, expiresAt)
	if err != nil {
		return fmt.Errorf("error while finding expired plan candidates: %v", err)
	}

	if len(*expiredPlanCandidateIds) == 0 {
		log.Println("No expired plan candidates found")
		return nil
	}

	log.Printf("Found %d expired plan candidates\n", len(*expiredPlanCandidateIds))

	// プラン候補を削除
	log.Printf("Deleting %d expired plan candidates\n", len(*expiredPlanCandidateIds))
	if err := s.planCandidateRepository.DeleteAll(ctx, *expiredPlanCandidateIds); err != nil {
		return fmt.Errorf("error while deleting expired plan candidates: %v", err)
	}
	log.Printf("Deleted %d expired plan candidates\n", len(*expiredPlanCandidateIds))

	return nil
}
