package plangen

import (
	"context"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
)

type CreatePlanParams struct {
	locationStart models.GeoLocation
	placeStart    models.Place
	places        []models.Place
}

// createPlanData 写真やタイトルなどのプランに必要な情報を作成する
func (s Service) createPlanData(ctx context.Context, planCandidateId string, params ...CreatePlanParams) []models.Plan {
	// レビュー・写真を取得する
	performanceTimer := time.Now()
	placeIdToPlaceWithPlaceDetail := s.fetchPlaceDetailData(ctx, planCandidateId, params...)
	s.logger.Info(
		"fetching reviews and images",
		zap.String("PlanCandidateId", planCandidateId),
		zap.Duration("duration", time.Since(performanceTimer)),
	)

	ch := make(chan *models.Plan, len(params))

	for _, param := range params {
		go func(ctx context.Context, param CreatePlanParams, ch chan<- *models.Plan) {
			// 出発地点から近い順に場所をめぐるように並び替え
			placesSortedByDistance := sortPlacesByDistanceFrom(param.locationStart, param.places)

			// プランのタイトルを生成
			chPlanTitle := make(chan string, 1)
			go func(ctx context.Context, chPlanTitle chan<- string) {
				performanceTimer := time.Now()
				title, err := s.GeneratePlanTitle(placesSortedByDistance)
				if err != nil {
					s.logger.Warn(
						"error while generating plan title",
						zap.String("PlanCandidateId", planCandidateId),
						zap.Error(err),
					)
					title = &param.placeStart.Google.Name
				}
				s.logger.Info(
					"generating plan title",
					zap.String("PlanCandidateId", planCandidateId),
					zap.Duration("duration", time.Since(performanceTimer)),
				)
				chPlanTitle <- *title
			}(ctx, chPlanTitle)

			// タイトル生成には2秒以上かかる場合があるため、タイムアウト処理を行う
			var title string
			chTitleTimeOut := time.NewTimer(2 * time.Second)
			select {
			case title = <-chPlanTitle:
				chTitleTimeOut.Stop()
			case <-chTitleTimeOut.C:
				s.logger.Warn(
					"timeout while generating plan title",
					zap.String("PlanCandidateId", planCandidateId),
				)
				title = param.placeStart.Google.Name
			}

			// プランに含まれる場所のレビューや写真をセットする
			var placesInPlan []models.Place
			for i, place := range param.places {
				if value, ok := placeIdToPlaceWithPlaceDetail[place.Id]; ok {
					param.places[i] = value
				}
				placesInPlan = append(placesInPlan, param.places[i])
			}

			ch <- &models.Plan{
				Id:     uuid.New().String(),
				Name:   title,
				Places: placesInPlan,
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

// fetchPlaceDetailData は、指定された場所の写真・レビューを一括で取得し、保存する
func (s Service) fetchPlaceDetailData(ctx context.Context, planCandidateId string, params ...CreatePlanParams) map[string]models.Place {
	// プラン間の場所の重複を無くすため、場所のIDをキーにして場所を保存する
	placeIdToPlace := make(map[string]models.Place)
	for _, param := range params {
		for _, place := range param.places {
			placeIdToPlace[place.Id] = place
		}

		// スタート地点（ユーザーが指定した場所 or スタート地点として選ばれた場所）も含める
		placeIdToPlace[param.placeStart.Id] = param.placeStart
	}

	// すべてのプランに含まれる Place を重複がないように選択し、写真を取得する
	placesToUpdate := make([]models.Place, 0, len(placeIdToPlace))
	for _, place := range placeIdToPlace {
		placesToUpdate = append(placesToUpdate, place)
	}

	placesToUpdate = s.placeSearchService.FetchPlacesPhotosAndSave(ctx, placesToUpdate...)

	for _, place := range placesToUpdate {
		placeIdToPlace[place.Id] = place
	}

	return placeIdToPlace
}
