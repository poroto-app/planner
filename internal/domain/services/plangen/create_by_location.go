package plangen

import (
	"context"
	"fmt"
	"googlemaps.github.io/maps"
	"log"
	"poroto.app/poroto/planner/internal/domain/array"
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
		placesSearched, err = s.placeService.SearchNearbyPlaces(ctx, locationStart)

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

	// レビューが低い、またはレビュー数が少ない場所を除外する
	placesFiltered = placefilter.FilterByRating(placesFiltered, 3.0, 10)

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

	// プラン作成の基準となる場所を選択
	var placesRecommend []places.Place

	if googlePlaceId != nil {
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

	placesRecommend = append(placesRecommend, s.selectBasePlace(
		placesFiltered,
		categoryNamesPreferred,
		categoryNamesDisliked,
		createBasedOnCurrentLocation,
	)...)
	for _, place := range placesRecommend {
		log.Printf("place recommended: %s\n", place.Name)
	}

	// 最もおすすめ度が高い３つの場所を基準にプランを作成する
	var createPlanParams []CreatePlanParams
	for _, placeRecommend := range placesRecommend {
		var placesInPlan []models.Place
		for _, createPlanParam := range createPlanParams {
			placesInPlan = append(placesInPlan, createPlanParam.places...)
		}

		planPlaces, err := s.createPlanPlaces(
			ctx,
			CreatePlanPlacesParams{
				locationStart:                locationStart,
				placeStart:                   placeRecommend,
				places:                       placesFiltered,
				placesOtherPlansContain:      placesInPlan,
				freeTime:                     freeTime,
				createBasedOnCurrentLocation: createBasedOnCurrentLocation,
				shouldOpenWhileTraveling:     createBasedOnCurrentLocation, // 現在地からプランを作成した場合は、今から出発した場合に閉まってしまうお店は含めない
			},
		)
		if err != nil {
			log.Printf("error while creating plan: %v\n", err)
			continue
		}

		createPlanParams = append(createPlanParams, CreatePlanParams{
			locationStart: locationStart,
			placeStart:    placeRecommend,
			places:        planPlaces,
		})
	}

	plans := s.createPlanData(ctx, createPlanSessionId, createPlanParams...)

	// 場所を指定してプランを作成した場合、その場所を起点としたプランを最初に表示する
	if googlePlaceId != nil {
		for i, plan := range plans {
			if len(plan.Places) == 0 {
				continue
			}

			firstPlace := plan.Places[0]
			if firstPlace.GooglePlaceId != nil && *firstPlace.GooglePlaceId == *googlePlaceId {
				plans[0], plans[i] = plans[i], plans[0]
				break
			}
		}
	}

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
