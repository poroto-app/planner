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
	expiredPlanCandidates, err := s.planCandidateRepository.FindExpiredBefore(ctx, expiresAt)
	if err != nil {
		return fmt.Errorf("error while finding expired plan candidates: %v", err)
	}

	log.Printf("Found %d expired plan candidates\n", len(*expiredPlanCandidates))

	planCandidateIds := make([]string, len(*expiredPlanCandidates))
	for i, expiredPlanCandidate := range *expiredPlanCandidates {
		planCandidateIds[i] = expiredPlanCandidate.Id
	}

	// 検索結果のキャッシュを削除
	log.Printf("Deleting %d expired place search results\n", len(planCandidateIds))
	if err := s.placeSearchResultRepository.DeleteAll(ctx, planCandidateIds); err != nil {
		return fmt.Errorf("error while deleting expired place search results: %v", err)
	}

	// プラン候補を削除
	log.Printf("Deleting %d expired plan candidates\n", len(planCandidateIds))
	if err := s.planCandidateRepository.DeleteAll(ctx, planCandidateIds); err != nil {
		return fmt.Errorf("error while deleting expired plan candidates: %v", err)
	}

	return nil
}
