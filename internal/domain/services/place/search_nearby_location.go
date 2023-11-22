package place

import (
	"context"
	"googlemaps.github.io/maps"
	"log"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/factory"
	"poroto.app/poroto/planner/internal/domain/models"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s Service) GetPlaceTypesToSearch() []maps.PlaceType {
	return []maps.PlaceType{
		"",
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

// GetPlaceTypesToPreSearch　ユーザの希望するカテゴリを質問するために検索する
func (s Service) GetPlaceTypesToPreSearch() []maps.PlaceType {
	return []maps.PlaceType{
		"",
		// カテゴリをしていなかった場合にヒットしないようなカテゴリを検索
		maps.PlaceTypeAmusementPark,
		maps.PlaceTypeShoppingMall,
		maps.PlaceTypeZoo,
	}
}

// GetPlaceTypesToDetailSearch プランを作成するために詳細に複数のカテゴリの場所を検索する
func (s Service) GetPlaceTypesToDetailSearch() []maps.PlaceType {
	placeTypes := make([]string, len(s.GetPlaceTypesToSearch()))
	for _, placeType := range s.GetPlaceTypesToSearch() {
		placeTypes = append(placeTypes, string(placeType))
	}

	// 事前に検索した場所で無いところを検索する
	var placeTypesToDetailSearch []maps.PlaceType
	for _, placeType := range s.GetPlaceTypesToSearch() {
		if !array.IsContain(placeTypes, string(placeType)) {
			placeTypesToDetailSearch = append(placeTypesToDetailSearch, placeType)
		}
	}

	return placeTypesToDetailSearch
}

// SearchNearbyPlaces location で指定された場所の付近にある場所を検索する
// また、特定のカテゴリに対して追加の検索を行う
func (s Service) SearchNearbyPlaces(ctx context.Context, location models.GeoLocation, placeTypes []maps.PlaceType) ([]models.GooglePlace, error) {
	placeTypesToSearch := s.GetPlaceTypesToPreSearch()

	ch := make(chan *[]models.GooglePlace, len(placeTypesToSearch))
	for _, placeType := range placeTypesToSearch {
		go func(ctx context.Context, ch chan<- *[]models.GooglePlace, placeType maps.PlaceType) {
			var placeTypePointer *maps.PlaceType
			if placeType != "" {
				placeTypePointer = &placeType
			}

			placesSearched, err := s.placesApi.FindPlacesFromLocation(ctx, &googleplaces.FindPlacesFromLocationRequest{
				Location: googleplaces.Location{
					Latitude:  location.Latitude,
					Longitude: location.Longitude,
				},
				Radius:      2000,
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
				places = append(places, factory.GooglePlaceFromPlaceEntity(place, nil, nil))
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

	return placesSearched, nil
}
