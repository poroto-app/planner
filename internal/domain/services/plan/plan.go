package plan

import (
	"context"
	"fmt"
	"log"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/domain/services/plangen"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type PlanService struct {
	placesApi                   places.PlacesApi
	planRepository              repository.PlanRepository
	planCandidateRepository     repository.PlanCandidateRepository
	placeSearchResultRepository repository.PlaceSearchResultRepository
	planGeneratorService        plangen.Service
}

func NewPlanService(ctx context.Context) (*PlanService, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	planRepository, err := firestore.NewPlanRepository(ctx)
	if err != nil {
		return nil, err
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	placeSearchResultRepository, err := firestore.NewPlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, err
	}

	planGeneratorService, err := plangen.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan generator service: %v", err)
	}

	return &PlanService{
		placesApi:                   *placesApi,
		planRepository:              planRepository,
		planCandidateRepository:     planCandidateRepository,
		placeSearchResultRepository: placeSearchResultRepository,
		planGeneratorService:        *planGeneratorService,
	}, err
}

func (s PlanService) CreatePlanByLocation(
	ctx context.Context,
	createPlanSessionId string,
	locationStart models.GeoLocation,
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

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	// 現在開店している場所だけを表示する
	placesFiltered = placefilter.FilterByOpeningNow(placesFiltered)

	// TODO: 移動距離ではなく、移動時間でやる
	var placesRecommend []places.Place

	placesInNear := placefilter.FilterWithinDistanceRange(placesFiltered, locationStart, 0, 500)
	placesInMiddle := placefilter.FilterWithinDistanceRange(placesFiltered, locationStart, 500, 1000)
	placesInFar := placefilter.FilterWithinDistanceRange(placesFiltered, locationStart, 1000, 2000)
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
			plan, err := s.planGeneratorService.CreatePlan(
				ctx,
				locationStart,
				placeRecommend,
				placesFiltered,
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

func (s PlanService) CreatePlanFromPlace(
	ctx context.Context,
	createPlanSessionId string,
	placeId string,
) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate")
	}

	// TODO: ユーザーの興味等を保存しておいて、それを反映させる
	placesSearched, err := s.placeSearchResultRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	var placeStart *places.Place
	for _, place := range placesSearched {
		if place.PlaceID == placeId {
			placeStart = &place
			break
		}
	}

	if placeStart == nil {
		return nil, fmt.Errorf("place not found")
	}

	planCreated, err := s.planGeneratorService.CreatePlan(
		ctx,
		placeStart.Location.ToGeoLocation(),
		*placeStart,
		placesSearched,
		// TODO: freeTimeの項目を保存し、それを反映させる
		nil,
		planCandidate.CreatedBasedOnCurrentLocation,
	)
	if err != nil {
		return nil, err
	}

	if _, err = s.planCandidateRepository.AddPlan(ctx, createPlanSessionId, planCreated); err != nil {
		return nil, err
	}

	return planCreated, nil
}
