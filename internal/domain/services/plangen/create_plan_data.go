package plangen

import (
	"context"
	"github.com/google/uuid"
	"log"
	"poroto.app/poroto/planner/internal/domain/models"
	api "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"time"
)

type CreatePlanParams struct {
	locationStart models.GeoLocation
	placeStart    api.Place
	places        []models.Place
}

// createPlanData 写真やタイトルなどのプランに必要な情報を作成する
func (s Service) createPlanData(ctx context.Context, planCandidateId string, params ...CreatePlanParams) []models.Plan {
	// 異なる複数のプランに所属するPlaceのIDの整合性を取る
	params = alignPlaceIds(params...)

	// 写真を取得する
	performanceTimer := time.Now()
	placeIdToImages := s.fetchAndSavePlacesPhotos(ctx, planCandidateId, params...)
	log.Printf("fetching place photos took %v\n", time.Since(performanceTimer))

	ch := make(chan *models.Plan, len(params))

	for _, param := range params {
		go func(ctx context.Context, param CreatePlanParams, ch chan<- *models.Plan) {
			places := param.places

			// プランのタイトルを生成
			chPlanTitle := make(chan string, 1)
			go func(ctx context.Context, places []models.Place, chPlanTitle chan<- string) {
				performanceTimer := time.Now()
				title, err := s.GeneratePlanTitle(param.places)
				if err != nil {
					log.Printf("error while generating plan title: %v\n", err)
					title = &param.placeStart.Name
				}
				log.Printf("generating plan title took %v\n", time.Since(performanceTimer))
				chPlanTitle <- *title
			}(ctx, places, chPlanTitle)

			// 場所のレビューを取得
			chPlansWithReviews := make(chan []models.Place, 1)
			go func(ctx context.Context, places []models.Place, chPlansWithReviews chan<- []models.Place) {
				performanceTimer := time.Now()
				places = s.FetchReviews(ctx, places)
				log.Printf("fetching place reviews took %v\n", time.Since(performanceTimer))
				chPlansWithReviews <- places
			}(ctx, places, chPlansWithReviews)

			// タイトル生成には2秒以上かかる場合があるため、タイムアウト処理を行う
			var title string
			chTitleTimeOut := time.NewTimer(2 * time.Second)
			select {
			case title = <-chPlanTitle:
				chTitleTimeOut.Stop()
			case <-chTitleTimeOut.C:
				log.Printf("timeout while generating plan title\n")
				title = param.placeStart.Name
			}

			placesWithReviews := <-chPlansWithReviews
			for i := 0; i < len(places); i++ {
				places[i].Images = placeIdToImages[places[i].Id]
				places[i].GooglePlaceReviews = placesWithReviews[i].GooglePlaceReviews
			}

			places = sortPlacesByDistanceFrom(param.locationStart, places)
			timeInPlan := planTimeFromPlaces(param.locationStart, places)

			ch <- &models.Plan{
				Id:            uuid.New().String(),
				Name:          title,
				Places:        places,
				TimeInMinutes: timeInPlan,
			}
		}(ctx, param, ch)
	}

	plans := make([]models.Plan, 0)
	for i := 0; i < len(params); i++ {
		plan := <-ch
		if plan == nil {
			continue
		}
		plans = append(plans, *plan)
	}

	return plans
}

// alignPlaceIds 複数のプランで出現する場所を同じIDにする
func alignPlaceIds(params ...CreatePlanParams) []CreatePlanParams {
	googlePlaceIdToPlaceId := make(map[string]string)
	for _, param := range params {
		for _, place := range param.places {
			if place.GooglePlaceId == nil {
				continue
			}

			googlePlaceIdToPlaceId[*place.GooglePlaceId] = place.Id
		}
	}

	for _, param := range params {
		for i, place := range param.places {
			if place.GooglePlaceId == nil {
				continue
			}

			param.places[i].Id = googlePlaceIdToPlaceId[*place.GooglePlaceId]
		}
	}

	return params
}

// fetchAndSavePlacesPhotos は、指定された場所の写真を一括で取得し、保存する
func (s Service) fetchAndSavePlacesPhotos(ctx context.Context, planCandidateId string, params ...CreatePlanParams) map[string][]models.Image {
	// プラン間の場所の重複を無くすため、場所のIDをキーにして場所を保存する
	placeIdToPlace := make(map[string]models.Place)
	for _, param := range params {
		for _, place := range param.places {
			placeIdToPlace[place.Id] = place
		}
	}

	// すべてのプランに含まれる Place を重複がないように選択し、写真を取得する
	places := make([]models.Place, 0, len(placeIdToPlace))
	for _, place := range placeIdToPlace {
		places = append(places, place)
	}
	places = s.FetchPlacesPhotos(ctx, places)

	// 取得した画像を保存
	for _, place := range places {
		if place.GooglePlaceId == nil {
			continue
		}

		if err := s.placeSearchResultRepository.SaveImagesIfNotExist(ctx, planCandidateId, *place.GooglePlaceId, place.Images); err != nil {
			log.Printf("error while saving images: %v\n", err)
			continue
		}
	}

	placeIdToImages := make(map[string][]models.Image)
	for _, place := range places {
		placeIdToImages[place.Id] = place.Images
	}

	return placeIdToImages
}
