package place

import (
	"context"
	"fmt"
	"googlemaps.github.io/maps"
	"log"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/factory"
	"poroto.app/poroto/planner/internal/domain/models"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// NearbySearchRadius NearbySearch で検索する際の半径
// FilterSearchResultRadius 検索結果をフィルタリングするときの半径
// IgnoreCategoryPlaceCount あるカテゴリの場所がこの数以上ある場合は、そのカテゴリの検索は行わない
const (
	NearbySearchRadius       = 2000
	FilterSearchResultRadius = 1000
	IgnoreCategoryPlaceCount = 5
)

type SearchNearbyPlacesInput struct {
	Location                 models.GeoLocation
	Radius                   uint
	FilterSearchResultRadius float64
	IgnoreCategoryPlaceCount uint
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

	if input.FilterSearchResultRadius == 0 {
		input.FilterSearchResultRadius = FilterSearchResultRadius
	}

	if input.IgnoreCategoryPlaceCount == 0 {
		input.IgnoreCategoryPlaceCount = IgnoreCategoryPlaceCount
	}

	// キャッシュされた検索結果を取得
	placesSaved, err := s.placeRepository.FindByLocation(ctx, input.Location)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places from location: %w", err)
	}

	// 検索箇所から半径 1000m 以内の場所を取得
	var placesFiltered []models.Place
	for _, place := range placesSaved {
		if place.Location.DistanceInMeter(input.Location) <= input.FilterSearchResultRadius {
			placesFiltered = append(placesFiltered, place)
		}
	}

	// 検索する必要のあるカテゴリを取得
	placeTypeToPlaces := groupByPlaceType(placesFiltered, s.placeTypesToSearch())
	var placeTypesToSearch []maps.PlaceType
	for placeType, places := range placeTypeToPlaces {
		// 5件以上の場所の検索結果が取得できた場合は、そのカテゴリの検索は行わない
		if len(places) >= int(input.IgnoreCategoryPlaceCount) {
			log.Printf("skip searching place type %s because it has enough places", placeType)
			continue
		}
		placeTypesToSearch = append(placeTypesToSearch, placeType)
	}

	ch := make(chan *[]models.GooglePlace, len(placeTypeToPlaces))
	for _, placeType := range placeTypesToSearch {
		go func(ctx context.Context, ch chan<- *[]models.GooglePlace, placeType maps.PlaceType) {
			var placeTypePointer *maps.PlaceType
			if placeType != "" {
				placeTypePointer = &placeType
			}

			placesSearched, err := s.placesApi.FindPlacesFromLocation(ctx, &googleplaces.FindPlacesFromLocationRequest{
				Location: googleplaces.Location{
					Latitude:  input.Location.Latitude,
					Longitude: input.Location.Longitude,
				},
				Radius:      input.Radius,
				Language:    "ja",
				Type:        placeTypePointer,
				SearchCount: 1,
			})
			if err != nil {
				ch <- nil
				log.Printf("error while fetching google_places with type %s: %v\n", placeType, err)
			}

			var places []models.GooglePlace
			for _, place := range placesSearched {
				places = append(places, factory.GooglePlaceFromPlaceEntity(place, nil))
			}

			ch <- &places
		}(ctx, ch, placeType)
	}

	// TODO：検索した場所の重複を削除する
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
	var placesSearchedFiltered []models.GooglePlace
	for _, place := range placesSearched {
		isAlreadyAdded := false
		for _, placeFiltered := range placesSearchedFiltered {
			if place.PlaceId == placeFiltered.PlaceId {
				isAlreadyAdded = true
				break
			}
		}

		if !isAlreadyAdded {
			placesSearchedFiltered = append(placesSearchedFiltered, place)
		}
	}

	return placesSearchedFiltered, nil
}

func (s Service) placeTypesToSearch() []maps.PlaceType {
	return []maps.PlaceType{
		maps.PlaceTypeAquarium,
		maps.PlaceTypeAmusementPark,
		maps.PlaceTypeCafe,
		maps.PlaceTypeMuseum,
		maps.PlaceTypeRestaurant,
		maps.PlaceTypeShoppingMall,
		maps.PlaceTypeSpa,
		maps.PlaceTypeZoo,
	}
}

func groupByPlaceType(places []models.Place, placeTypes []maps.PlaceType) map[maps.PlaceType][]models.Place {
	placesGroupedByPlaceType := make(map[maps.PlaceType][]models.Place)
	for _, placeType := range placeTypes {
		placesGroupedByPlaceType[placeType] = make([]models.Place, 0)

		for _, place := range places {
			if array.IsContain(place.Google.Types, string(placeType)) {
				placesGroupedByPlaceType[placeType] = append(placesGroupedByPlaceType[placeType], place)
			}
		}
	}

	return placesGroupedByPlaceType
}
