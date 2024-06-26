package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"time"
)

const (
	defaultMaxPlacesToRecommendPerCategory = 4
	defaultMaxPlacesToRecommend            = 16
)

type FetchPlacesToAddInput struct {
	PlanCandidateSetId string
	PlanId             string
	PlaceId            *string
	NLimit             uint
}

type FetchPlacesToAddOutput struct {
	PlacesRecommended []models.Place
	PlacesGrouped     []categoryGroupedPlaces
	PlacesAll         []models.Place
	Transitions       []models.Transition
}

type categoryGroupedPlaces struct {
	Category models.LocationCategory
	Places   []models.Place
}

// FetchPlacesToAdd はプランに追加する候補となる場所一覧を取得する
func (s Service) FetchPlacesToAdd(ctx context.Context, input FetchPlacesToAddInput) (*FetchPlacesToAddOutput, error) {
	if input.NLimit == 0 {
		input.NLimit = defaultMaxPlacesToRecommendPerCategory
	}

	if input.PlanCandidateSetId == "" {
		return nil, fmt.Errorf("plan candidate set id is empty")
	}

	if input.PlanId == "" {
		return nil, fmt.Errorf("plan id is empty")
	}

	planCandidateSet, err := s.planCandidateRepository.Find(ctx, input.PlanCandidateSetId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate set: %v", err)
	}

	var plan *models.Plan
	for _, p := range planCandidateSet.Plans {
		if p.Id == input.PlanId {
			plan = &p
			break
		}
	}
	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}

	if len(plan.Places) == 0 {
		return nil, fmt.Errorf("plan has no placesNearby")
	}

	var startPlace models.Place
	if input.PlaceId != nil {
		p, found := array.Find(plan.Places, func(place models.Place) bool {
			return place.Id == *input.PlaceId
		})
		if !found {
			return nil, fmt.Errorf("place(%s) not found in plan", *input.PlaceId)
		}
		startPlace = p
	} else {
		startPlace = plan.Places[0]
	}

	// 付近の場所を検索
	placesNearby, err := s.placeSearchService.SearchNearbyPlaces(ctx, placesearch.SearchNearbyPlacesInput{
		Location:           startPlace.Location,
		PlanCandidateSetId: &planCandidateSet.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching nearby places: %v\n", err)
	}

	categoriesToSearch := make([]models.LocationCategory, 0)

	// ユーザーが選択したカテゴリを優先的に調べる
	if planCandidateSet.MetaData.CategoriesPreferred != nil {
		categoriesToSearch = append(categoriesToSearch, *planCandidateSet.MetaData.CategoriesPreferred...)
	}

	for _, locationCategory := range []models.LocationCategory{
		models.CategoryRestaurant,
		models.CategoryCafe,
		models.CategoryShopping,
		models.CategoryAmusements,
		models.CategoryNatural,
		models.CategoryCulture,
		models.CategoryBakery,
		models.CategoryPark,
		models.CategorySpa,
	} {
		// すでに追加されている場合はスキップする
		_, isAlreadyContain := array.Find(categoriesToSearch, func(category models.LocationCategory) bool {
			return category.Name == locationCategory.Name
		})
		if isAlreadyContain {
			continue
		}

		// 検索対象から除外されている場合はスキップする
		if planCandidateSet.MetaData.CategoriesRejected != nil {
			_, isRejected := array.Find(*planCandidateSet.MetaData.CategoriesRejected, func(category models.LocationCategory) bool {
				return category.Name == locationCategory.Name
			})
			if isRejected {
				continue
			}
		}

		categoriesToSearch = append(categoriesToSearch, locationCategory)
	}

	// おすすめの場所を取得する
	placesRecommend := selectRecommendedPlaces(
		placesNearby,
		nil,
		*plan,
		startPlace.Location,
		planCandidateSet.MetaData,
		int(input.NLimit),
		nil,
	)

	// カテゴリごとのおすすめの場所を取得する
	var placesGrouped []categoryGroupedPlaces
	for _, category := range categoriesToSearch {
		placesAlreadyChosen := make([]models.Place, 0)
		placesAlreadyChosen = append(placesAlreadyChosen, placesRecommend...)
		placesAlreadyChosen = append(placesAlreadyChosen, array.FlatMap(placesGrouped, func(categoryGroupedPlaces categoryGroupedPlaces) []models.Place {
			return categoryGroupedPlaces.Places
		})...)
		placesAlreadyChosen = array.DistinctBy(placesAlreadyChosen, func(place models.Place) string {
			return place.Id
		})

		// 提案する場所の数が上限に達している場合はスキップする
		if len(placesRecommend) >= defaultMaxPlacesToRecommend {
			break
		}

		placesRecommendedWithCategory := selectRecommendedPlaces(
			placesNearby,
			placesAlreadyChosen,
			*plan,
			startPlace.Location,
			planCandidateSet.MetaData,
			int(input.NLimit),
			&category,
		)

		// ひとつも場所が見つからなかった場合はスキップする
		if len(placesRecommendedWithCategory) == 0 {
			continue
		}

		placesGrouped = append(placesGrouped, categoryGroupedPlaces{
			Category: category,
			Places:   placesRecommendedWithCategory,
		})
	}

	// 写真を取得
	var placesToFetchPhotos []models.Place
	placesToFetchPhotos = append(placesToFetchPhotos, placesRecommend...)
	placesToFetchPhotos = append(placesToFetchPhotos, array.FlatMap(placesGrouped, func(categoryGroupedPlaces categoryGroupedPlaces) []models.Place {
		return categoryGroupedPlaces.Places
	})...)
	placesToFetchPhotos = array.DistinctBy(placesToFetchPhotos, func(place models.Place) string {
		return place.Id
	})
	placesWithPhotos := s.placeSearchService.FetchPlacesPhotosAndSave(ctx, placesToFetchPhotos...)

	for i, place := range placesRecommend {
		placeWithPhoto, found := array.Find(placesWithPhotos, func(placeWithPhoto models.Place) bool {
			return placeWithPhoto.Id == place.Id
		})
		if !found {
			continue
		}
		placesRecommend[i] = placeWithPhoto
	}

	for iCategory, categoryGroupedPlaces := range placesGrouped {
		for iPlace, place := range categoryGroupedPlaces.Places {
			placeWithPhoto, found := array.Find(placesWithPhotos, func(placeWithPhoto models.Place) bool {
				return placeWithPhoto.Id == place.Id
			})
			if !found {
				continue
			}
			placesGrouped[iCategory].Places[iPlace] = placeWithPhoto
		}
	}

	// 移動時間の算出
	placesAll := placesRecommend
	placesAll = append(placesAll, array.FlatMap(placesGrouped, func(categoryGroupedPlaces categoryGroupedPlaces) []models.Place {
		return categoryGroupedPlaces.Places
	})...)
	placesAll = append(placesAll, startPlace)
	placesAll = array.DistinctBy(placesAll, func(place models.Place) string {
		return place.Id
	})
	transitions := array.Map(placesAll, func(place models.Place) models.Transition {
		return startPlace.CreateTransition(place)
	})

	return &FetchPlacesToAddOutput{
		PlacesRecommended: placesRecommend,
		PlacesGrouped:     placesGrouped,
		PlacesAll:         placesAll,
		Transitions:       transitions,
	}, nil
}

func selectRecommendedPlaces(
	places []models.Place,
	placesAlreadyChosen []models.Place,
	plan models.Plan,
	startLocation models.GeoLocation,
	planCandidateMetaData models.PlanCandidateMetaData,
	nLimit int,
	category *models.LocationCategory,
) []models.Place {
	placesFiltered := places
	placesFiltered = placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:        placesFiltered,
		StartLocation: startLocation,
	})

	// プランに含まれている場所から800m圏内の場所を選択する
	if len(plan.Places) > 1 {
		placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
			for _, placeInPlan := range plan.Places {
				if place.Location.DistanceInMeter(placeInPlan.Location) < 800 {
					return true
				}
			}
			return false
		})
	}

	// すでにプランに含まれている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
		_, isAlreadyInPlan := array.Find(plan.Places, func(placeInPlan models.Place) bool {
			return placeInPlan.Id == place.Id
		})
		return !isAlreadyInPlan
	})

	// すでに他のカテゴリで追加されている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
		_, isAlreadyChosen := array.Find(placesAlreadyChosen, func(placeAlreadyChosen models.Place) bool {
			return placeAlreadyChosen.Id == place.Id
		})
		return !isAlreadyChosen
	})

	// カテゴリでフィルタリング
	if category != nil {
		placesFiltered = placefilter.FilterByCategory(placesFiltered, []models.LocationCategory{*category}, true)
	}

	// 除外カテゴリでフィルタリング
	if planCandidateMetaData.CategoriesRejected != nil {
		placesFiltered = placefilter.FilterByCategory(placesFiltered, *planCandidateMetaData.CategoriesRejected, false)
	}

	// レビューの高い順でソート
	placesFiltered = models.SortPlacesByRating(placesFiltered)

	// TODO: 「カテゴリなし」の場合はすべてのカテゴリの場所が表示されるようにする
	placesRecommended := array.Take(placesFiltered, nLimit)

	return placesRecommended
}
