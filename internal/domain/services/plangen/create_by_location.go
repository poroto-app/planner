package plangen

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
	categoryNamesDisliked *[]string,
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
		placesSearched, err = s.SearchNearbyPlaces(ctx, locationStart)

		if err != nil {
			return nil, fmt.Errorf("error while fetching places: %v\n", err)
		}

		if err := s.placeSearchResultRepository.Save(ctx, createPlanSessionId, placesSearched); err != nil {
			log.Printf("error while saving places to cache: %v\n", err)
		}
		log.Printf("save places to cache[%v]\n", createPlanSessionId)
	}

	placesFiltered := placesSearched
	placesFiltered = placefilter.FilterIgnoreCategory(placesFiltered)
	placesFiltered = placefilter.FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)

	// 除外されたカテゴリがある場合はそのカテゴリを除外する
	if categoryNamesDisliked != nil {
		var categoriesDisliked []models.LocationCategory
		for _, categoryName := range *categoryNamesDisliked {
			if category := models.GetCategoryOfName(categoryName); category != nil {
				categoriesDisliked = append(categoriesDisliked, *category)
			}
		}
		placesFiltered = placefilter.FilterByCategory(placesFiltered, categoriesDisliked, false)
	}

	placesRecommend := s.selectBasePlace(
		placesFiltered,
		categoryNamesPreferred,
		categoryNamesDisliked,
		createBasedOnCurrentLocation,
	)
	for _, place := range placesRecommend {
		log.Printf("place recommended: %s\n", place.Name)
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
				placesFiltered,
				freeTime,
				createBasedOnCurrentLocation,
				// 現在地からプランを作成した場合は、今から出発した場合に閉まってしまうお店は含めない
				createBasedOnCurrentLocation,
			)
			if err != nil {
				log.Printf("error while creating plan: %v\n", err)
				chPlan <- nil
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
