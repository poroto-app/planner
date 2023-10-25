package plangen

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
)

type CreatePlanParams struct {
	locationStart models.GeoLocation
	placeStart    models.GooglePlace
	places        []models.GooglePlace
}

// createPlanData 写真やタイトルなどのプランに必要な情報を作成する
func (s Service) createPlanData(ctx context.Context, planCandidateId string, params ...CreatePlanParams) []models.Plan {
	// レビュー・写真・価格帯を取得する
	performanceTimer := time.Now()
	placeIdToPlaceDetailData := s.fetchPlaceDetailData(ctx, planCandidateId, params...)
	log.Printf("fetching reviews and images took %v\n", time.Since(performanceTimer))

	ch := make(chan *models.Plan, len(params))

	for _, param := range params {
		go func(ctx context.Context, param CreatePlanParams, ch chan<- *models.Plan) {
			googlePlaces := param.places

			googlePlaces = sortPlacesByDistanceFrom(param.locationStart, googlePlaces)
			timeInPlan := planTimeFromPlaces(param.locationStart, googlePlaces)

			// プランのタイトルを生成
			chPlanTitle := make(chan string, 1)
			go func(ctx context.Context, places []models.GooglePlace, chPlanTitle chan<- string) {
				performanceTimer := time.Now()
				title, err := s.GeneratePlanTitle(param.places)
				if err != nil {
					log.Printf("error while generating plan title: %v\n", err)
					title = &param.placeStart.Name
				}
				log.Printf("generating plan title took %v\n", time.Since(performanceTimer))
				chPlanTitle <- *title
			}(ctx, googlePlaces, chPlanTitle)

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

			var places []models.Place
			for i := 0; i < len(googlePlaces); i++ {
				if value, ok := placeIdToPlaceDetailData[googlePlaces[i].PlaceId]; ok {
					googlePlaces[i].Images = &value.Images
					googlePlaces[i].Reviews = &value.Reviews
					googlePlaces[i].PriceLevel = value.PriceLevel
				}
				places = append(places, googlePlaces[i].ToPlace())
			}

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

type placeDetail struct {
	GooglePlaceId string
	Reviews       []models.GooglePlaceReview
	Images        []models.Image
	PriceLevel    *int
}

// fetchPlaceDetailData は、指定された場所の写真・レビュー・値段帯を一括で取得し、保存する
func (s Service) fetchPlaceDetailData(ctx context.Context, planCandidateId string, params ...CreatePlanParams) map[string]placeDetail {
	// プラン間の場所の重複を無くすため、場所のIDをキーにして場所を保存する
	placeIdToPlace := make(map[string]models.GooglePlace)
	for _, param := range params {
		for _, place := range param.places {
			placeIdToPlace[place.PlaceId] = place
		}
	}

	// すべてのプランに含まれる Place を重複がないように選択し、写真を取得する
	places := make([]models.GooglePlace, 0, len(placeIdToPlace))
	for _, place := range placeIdToPlace {
		places = append(places, place)
	}

	places = s.placeService.FetchPlacesPhotosAndSave(ctx, planCandidateId, places...)
	places = s.placeService.FetchPlaceReviewsAndSave(ctx, planCandidateId, places...)
	places = s.placeService.FetchPlacesPriceLevelAndSave(ctx, planCandidateId, places...)

	placeIdToImages := make(map[string]placeDetail)
	for _, place := range places {
		var reviews []models.GooglePlaceReview
		var images []models.Image
		var priceLevel *int

		if place.Reviews != nil {
			reviews = *place.Reviews
		}

		if place.Images != nil {
			images = *place.Images
		}

		if place.PriceLevel != nil {
			priceLevel = place.PriceLevel
		}

		placeIdToImages[place.PlaceId] = placeDetail{
			GooglePlaceId: place.PlaceId,
			Reviews:       reviews,
			Images:        images,
			PriceLevel:    priceLevel,
		}
	}

	return placeIdToImages
}
