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
	placeStart    models.PlaceInPlanCandidate
	places        []models.PlaceInPlanCandidate
}

// createPlanData 写真やタイトルなどのプランに必要な情報を作成する
func (s Service) createPlanData(ctx context.Context, planCandidateId string, params ...CreatePlanParams) []models.Plan {
	// レビュー・写真を取得する
	performanceTimer := time.Now()
	placeIdToPlaceDetailData := s.fetchPlaceDetailData(ctx, planCandidateId, params...)
	log.Printf("fetching reviews and images took %v\n", time.Since(performanceTimer))

	ch := make(chan *models.Plan, len(params))

	for _, param := range params {
		go func(ctx context.Context, param CreatePlanParams, ch chan<- *models.Plan) {
			placesInPlanCandidate := param.places

			placesInPlanCandidate = sortPlacesByDistanceFrom(param.locationStart, placesInPlanCandidate)
			timeInPlan := planTimeFromPlaces(param.locationStart, placesInPlanCandidate)

			// プランのタイトルを生成
			chPlanTitle := make(chan string, 1)
			go func(ctx context.Context, chPlanTitle chan<- string) {
				performanceTimer := time.Now()
				title, err := s.GeneratePlanTitle(param.places)
				if err != nil {
					log.Printf("error while generating plan title: %v\n", err)
					title = &param.placeStart.Google.Name
				}
				log.Printf("generating plan title took %v\n", time.Since(performanceTimer))
				chPlanTitle <- *title
			}(ctx, chPlanTitle)

			// タイトル生成には2秒以上かかる場合があるため、タイムアウト処理を行う
			var title string
			chTitleTimeOut := time.NewTimer(2 * time.Second)
			select {
			case title = <-chPlanTitle:
				chTitleTimeOut.Stop()
			case <-chTitleTimeOut.C:
				log.Printf("timeout while generating plan title\n")
				title = param.placeStart.Google.Name
			}

			var places []models.Place
			for i := 0; i < len(placesInPlanCandidate); i++ {
				if value, ok := placeIdToPlaceDetailData[placesInPlanCandidate[i].Id]; ok {
					placesInPlanCandidate[i].Google.Images = &value.Images
					placesInPlanCandidate[i].Google.PlaceDetail = value.PlaceDetail
				}
				places = append(places, placesInPlanCandidate[i].ToPlace())
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
	PlaceId       string
	GooglePlaceId string
	Images        []models.Image
	PlaceDetail   *models.GooglePlaceDetail
}

// fetchPlaceDetailData は、指定された場所の写真・レビューを一括で取得し、保存する
func (s Service) fetchPlaceDetailData(ctx context.Context, planCandidateId string, params ...CreatePlanParams) map[string]placeDetail {
	// プラン間の場所の重複を無くすため、場所のIDをキーにして場所を保存する
	placeIdToPlace := make(map[string]models.PlaceInPlanCandidate)
	for _, param := range params {
		for _, place := range param.places {
			placeIdToPlace[place.Id] = place
		}
	}

	// すべてのプランに含まれる Place を重複がないように選択し、写真を取得する
	places := make([]models.PlaceInPlanCandidate, 0, len(placeIdToPlace))
	for _, place := range placeIdToPlace {
		places = append(places, place)
	}

	var googlePlaces []models.GooglePlace
	for _, place := range places {
		googlePlaces = append(googlePlaces, place.Google)
	}

	googlePlaces = s.placeService.FetchPlacesDetail(ctx, googlePlaces)
	googlePlaces = s.placeService.FetchPlacesPhotosAndSave(ctx, planCandidateId, googlePlaces...)

	placeIdToPlaceDetail := make(map[string]placeDetail)
	for _, place := range places {
		for _, googlePlace := range googlePlaces {
			if place.Google.PlaceId != googlePlace.PlaceId {
				continue
			}

			var images []models.Image

			if googlePlace.Images != nil {
				images = *googlePlace.Images
			}

			placeIdToPlaceDetail[place.Id] = placeDetail{
				PlaceId:       place.Id,
				GooglePlaceId: place.Google.PlaceId,
				Images:        images,
				PlaceDetail:   place.Google.PlaceDetail,
			}

			break
		}
	}

	return placeIdToPlaceDetail
}
