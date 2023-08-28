package plangen

import (
	"context"
	"googlemaps.github.io/maps"
	"log"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// SearchNearbyPlaces location で指定された場所の付近にある場所を検索する
// また、特定のカテゴリに対して追加の検索を行う
func (s Service) SearchNearbyPlaces(ctx context.Context, location models.GeoLocation) ([]places.Place, error) {
	var placeTypesToSearch = []maps.PlaceType{
		"",
		maps.PlaceTypeAquarium,
		maps.PlaceTypeAmusementPark,
		maps.PlaceTypeCafe,
		maps.PlaceTypeMovieTheater,
		maps.PlaceTypeMuseum,
		maps.PlaceTypeRestaurant,
		maps.PlaceTypeShoppingMall,
		maps.PlaceTypeSpa,
		maps.PlaceTypeZoo,
	}

	ch := make(chan *[]places.Place, len(placeTypesToSearch))
	for _, placeType := range placeTypesToSearch {
		go func(ctx context.Context, ch chan<- *[]places.Place, placeType maps.PlaceType) {
			var placeTypePointer *maps.PlaceType
			if placeType != "" {
				placeTypePointer = &placeType
			}

			placesSearched, err := s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
				Location: places.Location{
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
				log.Printf("error while fetching places with type %s: %v\n", placeType, err)
			}

			ch <- &placesSearched
		}(ctx, ch, placeType)
	}

	var placesSearched []places.Place
	for i := 0; i < len(placeTypesToSearch); i++ {
		searchResults := <-ch
		if searchResults == nil {
			continue
		}
		placesSearched = append(placesSearched, *searchResults...)
	}

	return placesSearched, nil
}
