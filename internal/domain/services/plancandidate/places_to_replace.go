package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

func (s *Service) FetchPlacesToReplace(
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

	if len(plan.Places) == 0 {
		return nil, fmt.Errorf("plan has no places")
	}

	startPlace := plan.Places[0]

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

	placesSearched, err := s.placeService.FetchSearchedPlaces(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching placesToReplace searched: %v", err)
	}

	placesFiltered := placesSearched

	// 遠い場所を除外
	placesFiltered = placefilter.FilterWithinDistanceRange(placesFiltered, startPlace.Location, 0, 1000)

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

	// 画像が取得できる場所のみを選択する
	placesFiltered = placefilter.FilterByHasPhoto(placesFiltered)

	// カテゴリのない場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
		return len(place.Categories()) > 0
	})

	// レビューの高い順でソート
	placesFiltered = models.SortPlacesByRating(placesFiltered)

	// 10件程度を候補として選択する
	if len(placesFiltered) > 10 {
		placesFiltered = placesFiltered[:10]
	}

	// なるべく一様に選択されるようにシャッフルする
	placesFiltered = models.ShufflePlaces(placesFiltered)

	// 指定された場所と同じカテゴリの場所が候補の半分の数だけ含められるようにする
	placesFilteredSameCategory := placefilter.FilterByCategory(placesFiltered, placeToReplace.Categories(), true)
	nPlacesSameCategory := int(nLimit / 2)
	if len(placesFilteredSameCategory) > nPlacesSameCategory {
		placesFilteredSameCategory = placesFilteredSameCategory[:nPlacesSameCategory]
	}

	// 指定された場所と異なるカテゴリの場所が候補の半分の数だけ含められるようにする
	placesFilteredDifferentCategory := placefilter.FilterByCategory(placesFiltered, placeToReplace.Categories(), false)
	nPlacesDifferentCategory := int(nLimit) - nPlacesSameCategory
	if len(placesFilteredDifferentCategory) > nPlacesDifferentCategory {
		placesFilteredDifferentCategory = placesFilteredDifferentCategory[:nPlacesDifferentCategory]
	}

	var placesToReplace []models.Place
	placesToReplace = append(placesToReplace, placesFilteredSameCategory...)
	placesToReplace = append(placesToReplace, placesFilteredDifferentCategory...)

	if len(placesToReplace) > int(nLimit) {
		placesToReplace = placesToReplace[:nLimit]
	}

	// 画像を取得
	placesToReplace = s.placeService.FetchPlacesPhotosAndSave(ctx, placesToReplace...)

	return placesToReplace, nil
}
