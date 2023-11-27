package main

import (
	"context"
	"log"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

func main() {
	env.LoadEnv()

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(context.Background())
	if err != nil {
		log.Fatalf("error while initializing plan candidate repository: %v", err)
	}

	// プラン候補を作成
	planCandidateId := "test"
	if err := planCandidateRepository.Create(context.Background(), planCandidateId, time.Now().Add(time.Hour*24*7)); err != nil {
		log.Fatalf("error while creating plan candidate: %v", err)
	}

	planCandidate, err := planCandidateRepository.Find(context.Background(), planCandidateId)
	if err != nil {
		log.Fatalf("error while finding plan candidate: %v", err)
	}
	log.Printf("plan candidate: %v", planCandidate)

	// 周囲の場所を検索後、その検索結果とプラン候補を結びつける
	searchedPlaceIds := []string{
		"0ZIll6UgUCgHLkrIv9Ze",
		"0WXayBy5zKBCmcAgWz30",
	}
	if err := planCandidateRepository.PlaceRepository.AddSearchedPlacesForPlanCandidate(context.Background(), planCandidateId, searchedPlaceIds); err != nil {
		log.Fatalf("error while adding searched places for plan candidate: %v", err)
	}

	// プランを保存
	if err := planCandidateRepository.AddPlan(context.Background(), planCandidateId, models.Plan{
		Id:   "test_plan",
		Name: "Text",
		Places: []models.Place{
			{
				Id: searchedPlaceIds[0],
			},
			{
				Id: searchedPlaceIds[1],
			},
		},
	}); err != nil {
		log.Fatalf("error while saving plan: %v", err)
	}

	// メタデータを保存
	if err := planCandidateRepository.UpdatePlanCandidateMetaData(context.Background(), planCandidateId, models.PlanCandidateMetaData{
		CategoriesPreferred: &[]models.LocationCategory{models.CategoryAmusements, models.CategoryBookStore},
		CategoriesRejected:  &[]models.LocationCategory{models.CategoryLibrary, models.CategoryCamp},
	}); err != nil {
		log.Fatalf("error while saving metadata for plan candidate: %v", err)
	}

	// 保存されたプラン候補を取得
	planCandidateSaved, err := planCandidateRepository.Find(context.Background(), planCandidateId)
	if err != nil {
		log.Fatalf("error while finding plan candidate: %v", err)
	}

	log.Printf("plan candidate: %v", planCandidateSaved)
}
