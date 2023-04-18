package places

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/array"
)

type LocationCategory struct {
	Name  string
	Photo PlacePhoto
}

func (r PlacesApi) FetchNearCategories(
	ctx context.Context,
	req *FindPlacesFromLocationRequest,
) ([]LocationCategory, error) {
	var nearCategories = []string{}
	var nearLocationCategories = []LocationCategory{}

	placesSearched, err := r.FindPlacesFromLocation(ctx, req)
	if err != nil {
		return nearLocationCategories, fmt.Errorf("error while fetching places: %v\n", err)
	}
	for _, place := range placesSearched {
		for _, category := range place.Categories {
			if !array.IsContain(nearCategories, category) {
				photos, err := r.FetchPlacePhotos(ctx, place)
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
						Photo: PlacePhoto{
							ImageUrl: "Not Found",
						},
					})
				}

			}
		}
	}
	return nearLocationCategories, nil
}
