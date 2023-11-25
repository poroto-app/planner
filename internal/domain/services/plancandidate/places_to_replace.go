package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"sort"
)

func (s Service) FetchPlacesToReplace(
	ctx context.Context,
	planCandidateId string,
	planId string,
	placeId string,
	nLimit uint,
) ([]models.Place, error) {
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

	var placeToReplace *models.Place
	for _, placeInPlan := range plan.Places {
		if placeInPlan.Id != placeId {
			placeToReplace = &placeInPlan
			break
		}
	}
	if placeToReplace == nil {
		return nil, fmt.Errorf("place to replace not found")
	}

	placesSearched, err := s.placeInPlanCandidateRepository.FindByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places searched: %v", err)
	}

	placesFiltered := *placesSearched

	// 重複した場所を削除
	placesFiltered = placefilter.FilterDuplicated(placesFiltered)

	// 会社はプランに含まれないようにする
	placesFiltered = placefilter.FilterCompany(placesFiltered)

	// 場所のカテゴリによるフィルタリング
	placesFiltered = placefilter.FilterIgnoreCategory(placesFiltered)
	placesFiltered = placefilter.FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)

	// すでにプランに含まれている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.PlaceInPlanCandidate) bool {
		for _, placeInPlan := range plan.Places {
			if placeInPlan.Id == place.Id {
				return false
			}
		}
		return true
	})

	// 指定された場所と同じカテゴリの場所を選択
	placesFiltered = placefilter.FilterByCategory(placesFiltered, placeToReplace.Categories, true)

	// レビューの高い順でソート
	sort.SliceStable(placesFiltered, func(i, j int) bool {
		return placesFiltered[i].Google.Rating > placesFiltered[j].Google.Rating
	})

	var places []models.PlaceInPlanCandidate
	for _, place := range placesFiltered {
		if len(place.Categories()) == 0 {
			continue
		}

		places = append(places, place)
	}

	if len(places) > int(nLimit) {
		places = places[:nLimit]
	}

	var googlePlacesToAdd []models.GooglePlace
	for _, place := range places {
		googlePlacesToAdd = append(googlePlacesToAdd, place.Google)
	}

	// 詳細情報を取得
	googlePlacesToAdd = s.placeService.FetchPlacesDetailAndSave(ctx, planCandidateId, googlePlacesToAdd)

	// 画像を取得
	googlePlacesToAdd = s.placeService.FetchPlacesPhotosAndSave(ctx, planCandidateId, googlePlacesToAdd...)

	var placesToAdd []models.Place
	for _, place := range places {
		for _, googlePlaceToAdd := range googlePlacesToAdd {
			if googlePlaceToAdd.PlaceId != place.Google.PlaceId {
				continue
			}

			place.Google = googlePlaceToAdd
			placesToAdd = append(placesToAdd, place.ToPlace())
		}
	}

	return placesToAdd, nil
}
