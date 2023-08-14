package plangenerator

import (
	"context"
	"fmt"
	"log"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s Service) CreatePlanByLocation(
	ctx context.Context,
	createPlanSessionId string,
	locationStart models.GeoLocation,
	// TODO: ユーザーに却下された場所を引数にする（プランを作成時により多くの場所を取得した場合、YESと答えたカテゴリの場所からしかプランを作成できなくなるため）
	categoryNamesPreferred *[]string,
	freeTime *int,
	createBasedOnCurrentLocation bool,
) (*[]models.Plan, error) {
	// 付近の場所を検索
	var placesSearched []places.Place

	//　キャッシュがあれば利用する
	placesCached, err := s.placeSearchResultRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		log.Printf("error while fetching places from cache: %v\n", err)
	} else if placesCached != nil {
		log.Printf("use cached places[%v]\n", createPlanSessionId)
		placesSearched = placesCached
	}

	if placesSearched == nil {
		placesSearched, err = s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
			Location: places.Location{
				Latitude:  locationStart.Latitude,
				Longitude: locationStart.Longitude,
			},
			Radius:   2000,
			Language: "ja",
		})

		if err != nil {
			return nil, fmt.Errorf("error while fetching places: %v\n", err)
		}

		if err := s.placeSearchResultRepository.Save(ctx, createPlanSessionId, placesSearched); err != nil {
			log.Printf("error while saving places to cache: %v\n", err)
		}
		log.Printf("save places to cache[%v]\n", createPlanSessionId)
	}

	var categoriesPreferred []models.LocationCategory
	if categoryNamesPreferred != nil {
		for _, categoryName := range *categoryNamesPreferred {
			if category := models.GetCategoryOfName(categoryName); category != nil {
				categoriesPreferred = append(categoriesPreferred, *category)
			}
		}
	}

	var categoriesToFiler []models.LocationCategory
	if len(*categoryNamesPreferred) > 0 {
		categoriesToFiler = categoriesPreferred
	} else {
		categoriesToFiler = models.GetCategoryToFilter()
	}

	placesFilter := placefilter.NewPlacesFilter(placesSearched)
	placesFilter = placesFilter.FilterIgnoreCategory()
	placesFilter = placesFilter.FilterByCategory(categoriesToFiler)

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesFilter = placesFilter.FilterByOpeningNow()

	// TODO: 移動距離ではなく、移動時間でやる
	var placesRecommend []places.Place

	placesInNear := placesFilter.FilterWithinDistanceRange(locationStart, 0, 500).Places()
	placesInMiddle := placesFilter.FilterWithinDistanceRange(locationStart, 500, 1000).Places()
	placesInFar := placesFilter.FilterWithinDistanceRange(locationStart, 1000, 2000).Places()
	if len(placesInNear) > 0 {
		// TODO: 0 ~ 500mで最もレビューの高い場所を選ぶ
		placesRecommend = append(placesRecommend, placesInNear[0])
	}
	if len(placesInMiddle) > 0 {
		// TODO: 500 ~ 1000mで最もレビューの高い場所を選ぶ
		placesRecommend = append(placesRecommend, placesInMiddle[0])
	}
	if len(placesInFar) > 0 {
		// TODO: 1000 ~ 2000mで最もレビューの高い場所を選ぶ
		placesRecommend = append(placesRecommend, placesInFar[0])
	}

	// 最もおすすめ度が高い３つの場所を基準にプランを作成する
	performanceTimer := time.Now()
	chPlans := make(chan *models.Plan, len(placesRecommend))
	for _, placeRecommend := range placesRecommend {
		go func(ctx context.Context, placeRecommend places.Place, chPlan chan<- *models.Plan) {
			plan, err := s.createPlan(
				ctx,
				locationStart,
				placeRecommend,
				placesFilter.Places(),
				freeTime,
				createBasedOnCurrentLocation,
			)
			if err != nil {
				log.Println(err)
				return
			}
			chPlans <- plan
		}(ctx, placeRecommend, chPlans)
	}

	plans := make([]models.Plan, 0)
	for i := 0; i < len(placesRecommend); i++ {
		plan := <-chPlans
		if plan == nil {
			continue
		}

		plans = append(plans, *plan)
	}
	log.Printf("created plans[%v]\n", time.Since(performanceTimer))

	return &plans, nil
}
