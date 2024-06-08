package batch

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
)

import (
	"context"
)

func DeleteExpiredPlanCandidateSet(ctx context.Context, db *sql.DB) error {
	log.Printf("=================== Start deleting expired plan candidates ===================\n")

	service, err := plancandidate.NewService(ctx, db)
	if err != nil {
		return fmt.Errorf("error while initializing plan candidate service: %v", err)
	}

	if err := service.DeleteExpiredPlanCandidates(ctx, time.Now()); err != nil {
		log.Printf("error while deleting expired plan candidates: %v", err)
		return fmt.Errorf("error while deleting expired plan candidates: %v", err)
	}

	log.Printf("=================== End deleting expired plan candidates ===================\n")

	return nil
}
