package plancandidate

import (
	"context"
	"fmt"
	"sort"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

// FetchPlacesToAdd はプランに追加する候補となる場所一覧を取得する
// nLimit によって取得する場所の数を制限することができる
func (s Service) FetchPlacesToAdd(ctx context.Context, planCandidateId string, planId string, nLimit uint) ([]models.Place, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}

	var plan *models.Plan
	for _, p := range planCandidate.Plans {
		if p.Id == planId {
			plan = &p
			break
		}
	}
	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}

	placesSearched, err := s.placeService.FetchSearchedPlaces(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places searched: %v", err)
	}

	placesFiltered := placesSearched

	// 重複した場所を削除
	placesFiltered = placefilter.FilterDuplicated(placesFiltered)

	// 会社はプランに含まれないようにする
	placesFiltered = placefilter.FilterCompany(placesFiltered)

	// 場所のカテゴリによるフィルタリング
	placesFiltered = placefilter.FilterIgnoreCategory(placesFiltered)
	placesFiltered = placefilter.FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)

	// すでにプランに含まれている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
		for _, placeInPlan := range plan.Places {
			if placeInPlan.Id == place.Id {
				return false
			}
		}
		return true
	})

	// レビューの高い順でソート
	sort.SliceStable(placesFiltered, func(i, j int) bool {
		return placesFiltered[i].Google.Rating > placesFiltered[j].Google.Rating
	})

	// TODO: すべてのカテゴリの場所が表示されるようにする
	var placesToAdd []models.Place
	for _, place := range placesFiltered {
		if len(placesToAdd) >= int(nLimit) {
			break
		}

		placesToAdd = append(placesToAdd, place)
	}

	// 場所の詳細情報を取得
	placesToAdd = s.placeService.FetchPlacesDetailAndSave(ctx, placesToAdd)

	// 写真を取得
	placesToAdd = s.placeService.FetchPlacesPhotosAndSave(ctx, placesToAdd...)

	return placesToAdd, nil
}
