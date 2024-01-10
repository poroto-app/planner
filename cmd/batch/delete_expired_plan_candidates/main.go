package main

import (
	"context"
	"log"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
	"time"

	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
)

func init() {
	env.LoadEnv()
}

// 有効期限切れのプラン候補を削除する
func main() {
	log.Printf("=================== Start deleting expired plan candidates ===================\n")

	db, err := rdb.InitDB(false)
	if err != nil {
		log.Fatalf("error while initializing db: %v", err)
	}

	service, err := plancandidate.NewService(db)
	if err != nil {
		log.Fatalf("error while initializing plan candidate service: %v", err)
	}

	ctx := context.Background()
	if err := service.DeleteExpiredPlanCandidates(ctx, time.Now()); err != nil {
		log.Fatalf("error while deleting expired plan candidates: %v", err)
	}

	log.Printf("=================== End deleting expired plan candidates ===================\n")
}
