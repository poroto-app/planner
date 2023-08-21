package plangen

import (
	"context"
	"fmt"
	"googlemaps.github.io/maps"
	"log"
	"poroto.app/poroto/planner/internal/domain/array"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s Service) CreatePlanByLocation(
	ctx context.Context,
	createPlanSessionId string,
	locationStart models.GeoLocation,
	// locationStart に対応する場所のID
	// これが指定されると、対応する場所を起点としてプランを作成する
	googlePlaceId *string,
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

	if googlePlaceId != nil {
		// TODO: 場所を指定された場合はプラン候補の最初に表示されるようにする
		// TODO: 他のplacesRecommendが指定された場所と近くならないようにする
		place, found, err := s.findOrFetchPlaceById(ctx, placesSearched, *googlePlaceId)
		if err != nil {
			log.Printf("error while fetching place: %v\n", err)
		}

		// TODO: キャッシュする

		// 開始地点となる場所が建物であれば、そこを基準としたプランを作成する
		if place != nil && array.IsContain(place.Types, string(maps.AutocompletePlaceTypeEstablishment)) {
			placesRecommend = append(placesRecommend, *place)
			if !found {
				placesFiltered = append(placesFiltered, *place)
			}
		}
	}

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
			plan, err := s.createPlan(
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

// findOrFetchPlaceById は、googlePlaceId に対応する場所を
// placesSearched から探し、なければAPIを使って取得する
func (s Service) findOrFetchPlaceById(
	ctx context.Context,
	placesSearched []places.Place,
	googlePlaceId string,
) (place *places.Place, found bool, err error) {
	for _, placeSearched := range placesSearched {
		if placeSearched.PlaceID == googlePlaceId {
			place = &placeSearched
			break
		}
	}

	if place != nil {
		return place, true, nil
	}

	place, err = s.placesApi.FetchPlace(ctx, places.FetchPlaceRequest{
		PlaceId:  googlePlaceId,
		Language: "ja",
	})
	if err != nil {
		return nil, false, fmt.Errorf("error while fetching place: %v\n", err)
	}

	return place, false, nil
}
