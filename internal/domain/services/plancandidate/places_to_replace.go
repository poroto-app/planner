package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
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

	placesSearched, err := s.placeSearchResultRepository.Find(ctx, planCandidateId)
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
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place places.Place) bool {
		for _, placeInPlan := range plan.Places {
			if placeInPlan.GooglePlaceId == nil {
				return false
			}

			if *placeInPlan.GooglePlaceId == place.PlaceID {
				return false
			}
		}
		return true
	})

	// 指定された場所と同じカテゴリの場所を選択
	placesFiltered = placefilter.FilterByCategory(placesFiltered, placeToReplace.Categories, true)

	// レビューの高い順でソート
	sort.SliceStable(placesFiltered, func(i, j int) bool {
		return placesFiltered[i].Rating > placesFiltered[j].Rating
	})

	var placesToAdd []models.Place
	for _, place := range placesFiltered {
		categories := models.GetCategoriesFromSubCategories(place.Types)
		if len(categories) == 0 {
			continue
		}

		// TODO: IDが統一されるようにする
		placesToAdd = append(placesToAdd, models.Place{
			Id:            place.PlaceID,
			Name:          place.Name,
			GooglePlaceId: utils.StrPointer(place.PlaceID),
			Location:      place.Location.ToGeoLocation(),
			Categories:    categories,
		})
	}

	placesToAdd = placesToAdd[:nLimit]

	// 写真を取得
	placesToAdd = s.placeService.FetchPlacesPhotosAndSave(ctx, planCandidateId, placesToAdd)

	// 口コミを取得
	placesToAdd = s.placeService.FetchPlaceReviewsAndSave(ctx, planCandidateId, placesToAdd)

	return placesToAdd, nil
}
