package plangen

import (
	"context"
	"fmt"
	"googlemaps.github.io/maps"
	"log"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
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
	var places []models.Place

	// すでに検索を行っている場合はその結果を取得
	placesSearched, err := s.placeService.FetchSearchedPlaces(ctx, createPlanSessionId)
	if err != nil {
		log.Printf("error while fetching searched places for plan candidate %v: %v\n", createPlanSessionId, err)
	} else if placesSearched != nil {
		log.Printf("use cached places[%v]\n", createPlanSessionId)
		places = placesSearched
	}

	// 検索を行っていない場合は検索を行う
	if places == nil {
		googlePlaces, err := s.placeService.SearchNearbyPlaces(ctx, place.SearchNearbyPlacesInput{Location: locationStart})
		if err != nil {
			return nil, fmt.Errorf("error while fetching google places: %v\n", err)
		}

		placesSaved, err := s.placeService.SaveSearchedPlaces(ctx, createPlanSessionId, googlePlaces)
		if err != nil {
			return nil, fmt.Errorf("error while saving searched places: %v\n", err)
		}

		places = placesSaved
	}

	placesFiltered := places
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

	log.Printf("places filtered: %v\n", len(placesFiltered))

	// プラン作成の基準となる場所を選択
	var placesRecommend []models.Place

	// 指定された場所の情報を取得する
	if googlePlaceId != nil {
		// TODO: 他のplacesRecommendが指定された場所と近くならないようにする
		place, found, err := s.findOrFetchPlaceById(ctx, createPlanSessionId, places, *googlePlaceId)
		if err != nil {
			log.Printf("error while fetching place: %v\n", err)
		}

		// 開始地点となる場所が建物であれば、そこを基準としたプランを作成する
		if place != nil && array.IsContain(place.Google.Types, string(maps.AutocompletePlaceTypeEstablishment)) {
			placesRecommend = append(placesRecommend, *place)
			if !found {
				placesFiltered = append(placesFiltered, *place)
			}
		}
	}

	// 場所を指定してプランを作成する場合は、指定した場所も含めて３つの場所を基準にプランを作成する
	maxBasePlaceCount := 3
	if googlePlaceId != nil {
		maxBasePlaceCount = 2
	}

	placesRecommend = append(placesRecommend, s.SelectBasePlace(SelectBasePlaceInput{
		Places:                 placesFiltered,
		CategoryNamesPreferred: categoryNamesPreferred,
		CategoryNamesDisliked:  categoryNamesDisliked,
		ShouldOpenNow:          false,
		MaxBasePlaceCount:      maxBasePlaceCount,
	})...)
	for _, place := range placesRecommend {
		log.Printf("place recommended: %s\n", place.Google.Name)
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
				planCandidateId:              createPlanSessionId,
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
			if firstPlace.Google.PlaceId == *googlePlaceId {
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
	planCandidateId string,
	placesSearched []models.Place,
	googlePlaceId string,
) (place *models.Place, found bool, err error) {
	for _, placeSearched := range placesSearched {
		if placeSearched.Google.PlaceId == googlePlaceId {
			place = &placeSearched
			break
		}
	}

	// すでに取得されている場合はそれを返す
	if place != nil {
		return place, true, nil
	}

	googlePlaceEntity, err := s.placeService.FetchGooglePlace(ctx, googlePlaceId)
	if err != nil {
		return nil, false, fmt.Errorf("error while fetching place: %v", err)
	}

	if googlePlaceEntity == nil {
		return nil, false, nil
	}

	// キャッシュする
	places, err := s.placeService.SaveSearchedPlaces(ctx, planCandidateId, []models.GooglePlace{*googlePlaceEntity})
	if err != nil {
		return nil, false, fmt.Errorf("error while saving searched places: %v", err)
	}

	if len(places) == 0 {
		return nil, false, fmt.Errorf("could not save searched places")
	}

	place = &places[0]

	return place, false, nil
}
