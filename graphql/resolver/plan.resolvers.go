package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.29

import (
	"context"
	"fmt"
	"log"

	"poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/services"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// CreatePlanByLocation is the resolver for the createPlanByLocation field.
func (r *mutationResolver) CreatePlanByLocation(ctx context.Context, input *model.CreatePlanByLocationInput) ([]*model.Plan, error) {
	// TODO: implement
	return []*model.Plan{}, nil
}

// MatchInterests is the resolver for the matchInterests field.
func (r *queryResolver) MatchInterests(ctx context.Context, input *model.MatchInterestsInput) (*model.InterestCandidate, error) {
	// TODO: 実際に付近の場所のカテゴリを提示する
	var categories = []*model.LocationCategory{}

	planService, err := services.NewPlanService()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}
	categoriesSearched, err := planService.FetchNearCategories(
		ctx,
		&places.FindPlacesFromLocationRequest{
			Location: places.Location{
				Latitude:  input.Latitude,
				Longitude: input.Longitude,
			},
			Radius: 2000,
		},
	)
	if err != nil {
		log.Println(err)

		log.Println(categoriesSearched)
	}

	for _, categorySearched := range categoriesSearched {
		categories = append(categories, &model.LocationCategory{
			Name:        categorySearched.Name,
			DisplayName: categorySearched.Name,
			Photo:       categorySearched.Photo.ImageUrl,
		})
	}
	return &model.InterestCandidate{
		Categories: categories,
	}, nil
}
