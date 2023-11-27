package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/env"
	repo "poroto.app/poroto/planner/internal/infrastructure/firestore"
)

const (
	planCandidateId = "test"
)

func init() {
	env.LoadEnv()
}

/**
 * このコードは、仮のプラン候補を用いて Firestore Plan Candidate Repository を利用するサンプルコードです。
 *
 * このコードを実行すると、以下の以下の動作が行われます。
 * 1. プラン候補を作成する
 * 2. プラン候補作成時に検索した場所のIDをプラン候補に紐づける
 * 3. プランを保存する
 * 4. メタデータを保存する
 * 5. 保存したプラン候補を取得する
 * 6. 保存したプラン候補の確認を行う
 * 7. 保存したプラン候補を削除する（CleanUp）
 **/
func main() {
	planCandidateRepository, err := repo.NewPlanCandidateRepository(context.Background())
	if err != nil {
		log.Fatalf("error while initializing plan candidate repository: %v", err)
	}

	// プラン候補を作成
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
	plan := models.Plan{
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
	}
	if err := planCandidateRepository.AddPlan(context.Background(), planCandidateId, plan); err != nil {
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

	// 保存されたプラン候補の確認
	if planCandidateSaved.Id != planCandidateId {
		log.Fatalf("plan candidate id is not matched: %s != %s", planCandidateSaved.Id, planCandidateId)
	}

	if planCandidateSaved.Plans[0].Id != plan.Id {
		log.Fatalf("plan id is not matched: %s != %s", planCandidateSaved.Plans[0].Id, "test_plan")
	}

	if planCandidateSaved.Plans[0].Places[0].Id != plan.Places[0].Id {
		log.Fatalf("place id is not matched: %s != %s", planCandidateSaved.Plans[0].Places[0].Id, plan.Places[0].Id)
	}

	// プラン候補を削除
	CleanUp(context.Background())
}

func CleanUp(ctx context.Context) error {
	planCandidateRepository, err := repo.NewPlanCandidateRepository(ctx)
	if err != nil {
		return fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	if err := planCandidateRepository.DeleteAll(ctx, []string{planCandidateId}); err != nil {
		return fmt.Errorf("error while deleting plan candidate: %v", err)
	}

	return nil
}
