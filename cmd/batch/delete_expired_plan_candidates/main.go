package main

import (
	"context"
	"log"
	"poroto.app/poroto/planner/internal/env"
	"time"

	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
)

func init() {
	env.LoadEnv()
}

// 有効期限切れのプラン候補を削除する
func main() {
	log.Printf("=================== Start deleting expired plan candidates ===================\n")

	ctx := context.Background()
	service, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Fatalf("error while initializing plan candidate service: %v", err)
	}

	if err := service.DeleteExpiredPlanCandidates(ctx, time.Now()); err != nil {
		log.Fatalf("error while deleting expired plan candidates: %v", err)
	}

	log.Printf("=================== End deleting expired plan candidates ===================\n")
}
