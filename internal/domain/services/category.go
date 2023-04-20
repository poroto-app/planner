package services

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type LocationCategory struct {
	Name  string
	Photo places.PlacePhoto
}

func (s PlanService) FetchNearCategories(
	ctx context.Context,
	req *places.FindPlacesFromLocationRequest,
) ([]LocationCategory, error) {
	var nearCategories = []string{}
	var nearLocationCategories = []LocationCategory{}

	placesSearched, err := s.placesApi.FindPlacesFromLocation(ctx, req)
	if err != nil {
		return nearLocationCategories, fmt.Errorf("error while fetching places: %v\n", err)
	}
	for _, place := range placesSearched {
		for _, category := range place.Types {
			if !array.IsContain(nearCategories, category) {
				photos, err := s.placesApi.FetchPlacePhotos(ctx, place)
				if err != nil {
					continue
				}
				nearCategories = append(nearCategories, category)

				if photos != nil {
					nearLocationCategories = append(nearLocationCategories, LocationCategory{
						Name:  category,
						Photo: photos[0],
					})
				} else {
					nearLocationCategories = append(nearLocationCategories, LocationCategory{
						Name: category,
						Photo: places.PlacePhoto{
							ImageUrl: "Not Found",
						},
					})
				}

			}
		}
	}
	return nearLocationCategories, nil
}
