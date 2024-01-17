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
)

// NearbySearchRadius NearbySearch で検索する際の半径
// FilterSearchResultRadius 検索結果をフィルタリングするときの半径
const (
	NearbySearchRadius = 5 * 1000
)

type SearchNearbyPlacesInput struct {
	Location models.GeoLocation
	Radius   uint
}

// placeTypeWithCondition 検索する必要のあるカテゴリを表す
// searchRange Nearby Search時の検索範囲（水族館等の施設の数が少ない場所を探すときは広い範囲を探す）
// filterRange 周囲ににあるかどうかを確認するときの検索範囲 (カフェ等の施設の数が多い場所を探すときは狭い範囲を探す)
// ignorePlaceCount あるカテゴリの場所がこの数以上ある場合は、そのカテゴリの検索は行わない
type placeTypeWithCondition struct {
	placeType        maps.PlaceType
	searchRange      uint
	filterRange      uint
	ignorePlaceCount uint
}

// SearchNearbyPlaces location で指定された場所の付近にある場所を検索する
// また、特定のカテゴリに対して追加の検索を行う
func (s Service) SearchNearbyPlaces(ctx context.Context, input SearchNearbyPlacesInput) ([]models.GooglePlace, error) {
	if input.Location.Latitude == 0 || input.Location.Longitude == 0 {
		panic("location is not specified")
	}

	if input.Radius == 0 {
		input.Radius = NearbySearchRadius
	}

	// キャッシュされた検索結果を取得
	placesSaved, err := s.placeRepository.FindByLocation(ctx, input.Location, float64(input.Radius))
	if err != nil {
		return nil, fmt.Errorf("error while fetching places from location: %w", err)
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
	placeTypeToPlaces := groupByPlaceType(placesSaved, s.placeTypesToSearch())
	var placeTypesToSearch []placeTypeWithCondition
	for _, placeTypeToSearch := range s.placeTypesToSearch() {
		savedPlacesOfPlaceType := placeTypeToPlaces[placeTypeToSearch.placeType]

		// 保存された場所の中から特定の範囲内にある場所を取得
		placesInSearchRange := placefilter.FilterWithinDistanceRange(savedPlacesOfPlaceType, input.Location, 0, float64(placeTypeToSearch.filterRange))

		s.logger.Info(
			"saved places of place type",
			zap.String("placeType", string(placeTypeToSearch.placeType)),
			zap.Int("placesInSearchRange", len(placesInSearchRange)),
			zap.Uint("ignorePlaceCount", placeTypeToSearch.ignorePlaceCount),
		)

		// 必要な分だけ場所の検索結果が取得できた場合は、そのカテゴリの検索は行わない
		if len(placesInSearchRange) >= int(placeTypeToSearch.ignorePlaceCount) {
			s.logger.Debug(
				"skip searching place type because it has enough places",
				zap.String("placeType", string(placeTypeToSearch.placeType)),
				zap.Int("savedPlacesOfPlaceType", len(savedPlacesOfPlaceType)),
				zap.Int("placesInSearchRange", len(placesInSearchRange)),
				zap.Uint("ignorePlaceCount", placeTypeToSearch.ignorePlaceCount),
			)
			continue
		}

		placeTypesToSearch = append(placeTypesToSearch, placeTypeToSearch)
	}

	ch := make(chan *[]models.GooglePlace, len(placeTypeToPlaces))
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

	var placesSearched []models.GooglePlace
	for i := 0; i < len(placeTypesToSearch); i++ {
		searchResults := <-ch
		if searchResults == nil {
			continue
		}
		placesSearched = append(placesSearched, *searchResults...)
	}

	// 検索された場所に加えて、キャッシュされた場所を追加
	for _, place := range placesSaved {
		placesSearched = append(placesSearched, place.Google)
	}

	// 重複した場所を削除
	placesSearchedFiltered := array.DistinctBy(placesSearched, func(place models.GooglePlace) string {
		return place.PlaceId
	})

	return placesSearchedFiltered, nil
}

func (s Service) placeTypesToSearch() []placeTypeWithCondition {
	// そのカテゴリの場所が filterRange で指定している範囲の中に
	// このくらいはありそうという値を ignorePlaceCount に指定している
	// また、検索できる最大サイズは50kmまで
	return []placeTypeWithCondition{
		// 付近になければ、一度も検索していないことを怪しむレベル
		{
			placeType:        maps.PlaceTypeCafe,
			searchRange:      10 * 1000,
			filterRange:      5 * 1000,
			ignorePlaceCount: 3,
		},
		{
			placeType:        maps.PlaceTypeRestaurant,
			searchRange:      10 * 1000,
			filterRange:      5 * 1000,
			ignorePlaceCount: 5,
		},
		// 近くにあってもおかしくないレベル
		{
			placeType:        maps.PlaceTypeShoppingMall,
			searchRange:      20 * 1000,
			filterRange:      5 * 1000,
			ignorePlaceCount: 1,
		},
		{
			placeType:        maps.PlaceTypeSpa,
			searchRange:      50 * 1000,
			filterRange:      10 * 1000,
			ignorePlaceCount: 1,
		},
		{
			placeType:        maps.PlaceTypeTouristAttraction,
			searchRange:      30 * 1000,
			filterRange:      10 * 1000,
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
			searchRange:      50 * 1000,
			filterRange:      30 * 1000,
			ignorePlaceCount: 1,
		},
		{
			placeType:        maps.PlaceTypeMuseum,
			searchRange:      50 * 1000,
			filterRange:      30 * 1000,
			ignorePlaceCount: 1,
		},
		{
			placeType:        maps.PlaceTypeZoo,
			searchRange:      50 * 1000,
			filterRange:      30 * 1000,
			ignorePlaceCount: 1,
		},
	}
}

func groupByPlaceType(places []models.Place, placeTypes []placeTypeWithCondition) map[maps.PlaceType][]models.Place {
	placesGroupedByPlaceType := make(map[maps.PlaceType][]models.Place)
	for _, placeType := range placeTypes {
		placesGroupedByPlaceType[placeType.placeType] = make([]models.Place, 0)

		for _, place := range places {
			if array.IsContain(place.Google.Types, string(placeType.placeType)) {
				placesGroupedByPlaceType[placeType.placeType] = append(placesGroupedByPlaceType[placeType.placeType], place)
			}
		}
	}

	return placesGroupedByPlaceType
}
