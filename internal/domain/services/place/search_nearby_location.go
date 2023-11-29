package place

import (
	"context"
	"googlemaps.github.io/maps"
	"log"
	"poroto.app/poroto/planner/internal/domain/factory"
	"poroto.app/poroto/planner/internal/domain/models"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// SearchNearbyPlaces location で指定された場所の付近にある場所を検索する
// また、特定のカテゴリに対して追加の検索を行う
func (s Service) SearchNearbyPlaces(ctx context.Context, location models.GeoLocation) ([]models.GooglePlace, error) {
	var placeTypesToSearch = []maps.PlaceType{
		//"",
		//maps.PlaceTypeAquarium,
		//maps.PlaceTypeAmusementPark,
		maps.PlaceTypeCafe,
		//maps.PlaceTypeMuseum,
		//maps.PlaceTypeRestaurant,
		//maps.PlaceTypeShoppingMall,
		//maps.PlaceTypeSpa,
		//maps.PlaceTypeZoo,
	}

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
