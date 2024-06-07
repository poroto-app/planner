package placesearch

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/factory"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"time"
)

// nearbySearchRadius すでに保存された場所から近くにある場所を検索するときの検索範囲
const (
	nearbySearchRadius = 5 * 1000
)

// SearchNearbyPlacesInput は付近の場所を検索するときの入力
//
// PlanCandidateSetId が指定されており、すでに対応するプラン候補で検索が行われている場合は、検索を行わない
// これは、付近の特定のカテゴリが無い場合に、すでに検索が行われているのに、再度検索が行われてしまうという状況を防ぐため
// 例：水族館は30km圏内に存在しない場合は検索が行われるが、検索したことを記録していないと、水族館が30km圏内に存在しない場合にもう一度検索が行われてしまう
type SearchNearbyPlacesInput struct {
	Location           models.GeoLocation
	PlanCandidateSetId *string
}

// placeTypeWithCondition 検索する必要のあるカテゴリを表す
//
// searchRange Nearby Search時の検索範囲（水族館等の施設の数が少ない場所を探すときは広い範囲を探す）
// この距離を指定すればぴったり20件取得できるという値を指定している(検索できる最大サイズは50kmまで)
//
// filterRange 周囲にあるかどうかを確認するときの検索範囲 (カフェ等の施設の数が多い場所を探すときは狭い範囲を探す)
// 最低限この範囲にはあるはずという値を指定している
//
// ignorePlaceCount あるカテゴリの場所がこの数以上ある場合は、そのカテゴリの検索は行わない
// このくらいはありそうという値を指定している
type placeTypeWithCondition struct {
	placeType        maps.PlaceType
	searchRange      uint
	filterRange      uint
	ignorePlaceCount uint
}

// SearchNearbyPlaces location で指定された場所の付近にある場所を検索し、保存する
// また、特定のカテゴリに対して追加の検索を行う
func (s Service) SearchNearbyPlaces(ctx context.Context, input SearchNearbyPlacesInput) ([]models.Place, error) {
	if input.Location.Latitude == 0 || input.Location.Longitude == 0 {
		panic("location is not specified")
	}

	var isAlreadySearched bool
	if input.PlanCandidateSetId != nil {
		planCandidateSet, err := s.planCandidateRepository.Find(ctx, *input.PlanCandidateSetId, time.Now())
		if err != nil {
			// エラーが発生しても処理を続行する
			s.logger.Warn("error while fetching plan candidate set", zap.Error(err))
		} else if planCandidateSet != nil {
			isAlreadySearched = planCandidateSet.IsPlaceSearched
		}
	} else {
		s.logger.Warn("plan candidate set id is not specified")
	}

	// キャッシュされた検索結果を取得
	placesSaved, err := s.placeRepository.FindByLocation(ctx, input.Location, float64(nearbySearchRadius))
	if err != nil {
		return nil, fmt.Errorf("error while fetching places from location: %w", err)
	}

	// すでに検索されている場合は、検索を行わない
	if isAlreadySearched && len(placesSaved) > 0 {
		s.logger.Info("skip searching places because it has already been searched", zap.Int("places", len(placesSaved)))
		return placesSaved, nil
	}

	// カテゴリごとにキャッシュされた検索結果を取得
	for _, placeTypeWithCondition := range s.placeTypesToSearch() {
		placesSavedWithType, err := s.placeRepository.FindByGooglePlaceType(
			ctx,
			string(placeTypeWithCondition.placeType),
			input.Location,
			float64(placeTypeWithCondition.searchRange),
		)
		if err != nil {
			return nil, fmt.Errorf("error while fetching places from google place type: %w", err)
		}
		placesSaved = append(placesSaved, *placesSavedWithType...)
	}

	// 重複した場所を削除
	placesSaved = array.DistinctBy(placesSaved, func(place models.Place) string { return place.Id })
	s.logger.Info("successfully fetched saved places", zap.Int("places", len(placesSaved)))

	// 検索する必要のあるカテゴリを取得
	var placeTypesToSearch []placeTypeWithCondition
	for _, placeTypeToSearch := range s.placeTypesToSearch() {
		savedPlacesOfPlaceType := array.Filter(placesSaved, func(place models.Place) bool {
			if placeTypeToSearch.placeType == "" {
				return true
			}
			return array.IsContain(place.Google.Types, string(placeTypeToSearch.placeType))
		})

		// 保存された場所の中から特定の範囲内にある場所を取得
		placesOfPlaceTypeInRange := placefilter.FilterWithinDistanceRange(savedPlacesOfPlaceType, input.Location, 0, float64(placeTypeToSearch.filterRange))

		s.logger.Info(
			"saved places of place type",
			zap.String("placeType", string(placeTypeToSearch.placeType)),
			zap.Int("placesOfPlaceTypeInRange", len(placesOfPlaceTypeInRange)),
			zap.Uint("ignorePlaceCount", placeTypeToSearch.ignorePlaceCount),
		)

		// 必要な分だけ場所の検索結果が取得できた場合は、そのカテゴリの検索は行わない
		if len(placesOfPlaceTypeInRange) >= int(placeTypeToSearch.ignorePlaceCount) {
			s.logger.Debug(
				"skip searching place type because it has enough places",
				zap.String("placeType", string(placeTypeToSearch.placeType)),
				zap.Int("savedPlacesOfPlaceType", len(savedPlacesOfPlaceType)),
				zap.Int("placesOfPlaceTypeInRange", len(placesOfPlaceTypeInRange)),
				zap.Uint("ignorePlaceCount", placeTypeToSearch.ignorePlaceCount),
			)
			continue
		}

		placeTypesToSearch = append(placeTypesToSearch, placeTypeToSearch)
	}

	ch := make(chan *[]models.GooglePlace, len(placeTypesToSearch))
	for _, placeType := range placeTypesToSearch {
		go func(ctx context.Context, ch chan<- *[]models.GooglePlace, placeTypeWithCondition placeTypeWithCondition) {
			var placeTypePointer *maps.PlaceType
			if placeTypeWithCondition.placeType != "" {
				placeTypePointer = &placeTypeWithCondition.placeType
			}

			placesSearched, err := s.placesApi.NearbySearch(ctx, &googleplaces.NearbySearchRequest{
				Location: googleplaces.Location{
					Latitude:  input.Location.Latitude,
					Longitude: input.Location.Longitude,
				},
				Radius:      placeTypeWithCondition.searchRange,
				Language:    "ja",
				Type:        placeTypePointer,
				SearchCount: 1,
			})
			if err != nil {
				// TODO: channelを用いてエラーハンドリングする
				ch <- nil
				s.logger.Warn(
					"error while fetching google_places",
					zap.String("placeType", string(placeTypeWithCondition.placeType)),
					zap.Uint("searchRange", placeTypeWithCondition.searchRange),
					zap.Error(err),
				)
			}

			s.logger.Info(
				"successfully fetched nearby places",
				zap.String("placeType", string(placeTypeWithCondition.placeType)),
				zap.Uint("searchRange", placeTypeWithCondition.searchRange),
				zap.Int("places", len(placesSearched)),
			)

			var places []models.GooglePlace
			for _, place := range placesSearched {
				places = append(places, factory.GooglePlaceFromPlaceEntity(place, nil))
			}

			ch <- &places
		}(ctx, ch, placeType)
	}

	var googlePlacesSearched []models.GooglePlace
	for i := 0; i < len(placeTypesToSearch); i++ {
		searchResults := <-ch
		if searchResults == nil {
			continue
		}
		googlePlacesSearched = append(googlePlacesSearched, *searchResults...)
	}

	// 検索された場所に加えて、キャッシュされた場所を追加
	for _, place := range placesSaved {
		googlePlacesSearched = append(googlePlacesSearched, place.Google)
	}

	// 重複した場所を削除
	placesSearchedFiltered := array.DistinctBy(googlePlacesSearched, func(place models.GooglePlace) string {
		return place.PlaceId
	})

	// 検索された場所を保存
	places, err := s.placeRepository.SavePlacesFromGooglePlaces(ctx, placesSearchedFiltered...)
	if err != nil {
		return nil, fmt.Errorf("error while saving places from google place: %v\n", err)
	}

	// 検索済みであることを保存する
	if input.PlanCandidateSetId != nil {
		err := s.planCandidateRepository.UpdateIsPlaceSearched(ctx, *input.PlanCandidateSetId, true)
		if err != nil {
			s.logger.Warn("error while updating place searched", zap.Error(err))
		}
	}

	if places == nil {
		return nil, nil
	}

	return *places, nil
}

func (s Service) placeTypesToSearch() []placeTypeWithCondition {
	return []placeTypeWithCondition{
		// 付近になければ、一度も検索していないことを怪しむレベル
		{
			placeType:        "",
			searchRange:      2 * 1000,
			filterRange:      1 * 1000,
			ignorePlaceCount: 5,
		},
		{
			placeType:        maps.PlaceTypeCafe,
			searchRange:      3 * 1000,
			filterRange:      3 * 1000,
			ignorePlaceCount: 5,
		},
		{
			placeType:        maps.PlaceTypeBookStore,
			searchRange:      3 * 1000,
			filterRange:      3 * 1000,
			ignorePlaceCount: 5,
		},
		{
			placeType:        maps.PlaceTypeRestaurant,
			searchRange:      3 * 1000,
			filterRange:      3 * 1000,
			ignorePlaceCount: 5,
		},
		{
			placeType:        maps.PlaceTypeBakery,
			searchRange:      3 * 1000,
			filterRange:      5 * 1000,
			ignorePlaceCount: 3,
		},
		{
			placeType:        maps.PlaceTypeTouristAttraction,
			searchRange:      3 * 1000,
			filterRange:      1 * 1000,
			ignorePlaceCount: 1,
		},
		// 近くにあってもおかしくないレベル
		{
			placeType:        maps.PlaceTypeShoppingMall,
			searchRange:      5 * 1000,
			filterRange:      5 * 1000,
			ignorePlaceCount: 1,
		},
		{
			placeType:        maps.PlaceTypeSpa,
			searchRange:      5 * 1000,
			filterRange:      5 * 1000,
			ignorePlaceCount: 1,
		},
		{
			placeType:        maps.PlaceTypeMuseum,
			searchRange:      10 * 1000,
			filterRange:      5 * 1000,
			ignorePlaceCount: 1,
		},
		// 近くに無いことがあたりまえなレベル
		{
			placeType:        maps.PlaceTypeAquarium,
			searchRange:      50 * 1000,
			filterRange:      30 * 1000,
			ignorePlaceCount: 1,
		},
		{
			placeType:        maps.PlaceTypeAmusementPark,
			searchRange:      20 * 1000,
			filterRange:      10 * 1000,
			ignorePlaceCount: 1,
		},
		{
			placeType:        maps.PlaceTypeZoo,
			searchRange:      30 * 1000,
			filterRange:      10 * 1000,
			ignorePlaceCount: 1,
		},
	}
}
