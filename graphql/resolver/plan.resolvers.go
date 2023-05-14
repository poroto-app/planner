package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.29

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services"
)

// CreatePlanByLocation is the resolver for the createPlanByLocation field.
func (r *mutationResolver) CreatePlanByLocation(ctx context.Context, input model.CreatePlanByLocationInput) (*model.CreatePlanByLocationOutput, error) {
	// TODO: エラー時は処理を停止させる
	service, err := services.NewPlanService()
	if err != nil {
		log.Println(err)
	}

	plans, err := service.CreatePlanByLocation(
		ctx,
		models.GeoLocation{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		})
	if err != nil {
		log.Println(err)
	}

	return &model.CreatePlanByLocationOutput{
		Session: uuid.New().String(),
		Plans:   plansFromDomainModel(plans),
	}, nil
}

// MatchInterests is the resolver for the matchInterests field.
func (r *queryResolver) MatchInterests(ctx context.Context, input *model.MatchInterestsInput) (*model.InterestCandidate, error) {
	planService, err := services.NewPlanService()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	categoriesSearched, err := planService.CategoriesNearLocation(
		ctx,
		models.GeoLocation{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while searching categories: %v", err)
	}

	var categories = []*model.LocationCategory{}
	for _, categorySearched := range categoriesSearched {
		categories = append(categories, &model.LocationCategory{
			Name:        categorySearched.Name,
			DisplayName: categorySearched.DisplayName,
			Photo:       categorySearched.Photo,
		})
	}
	return &model.InterestCandidate{
		Categories: categories,
	}, nil
}

// CachedCreatedPlans is the resolver for the CachedCreatedPlans field.
func (r *queryResolver) CachedCreatedPlans(ctx context.Context, input model.CachedCreatedPlansInput) (*model.CachedCreatedPlans, error) {
	return &model.CachedCreatedPlans{
		Plans: nil,
	}, nil
}

func plansFromDomainModel(plans *[]models.Plan) []*model.Plan {
	graphqlPlans := make([]*model.Plan, 0)

	for _, plan := range *plans {
		places := make([]*model.Place, 0)
		for _, place := range plan.Places {
			places = append(places, &model.Place{
				Name:   place.Name,
				Photos: place.Photos,
				Location: &model.GeoLocation{
					Latitude:  place.Location.Latitude,
					Longitude: place.Location.Longitude,
				},
				EstimatedStayDuration: int(place.EstimatedStayDuration),
			})
		}

		graphqlPlans = append(graphqlPlans, &model.Plan{
			ID:            plan.Id,
			Name:          plan.Name,
			Places:        places,
			TimeInMinutes: int(plan.TimeInMinutes),
		})
	}

	return graphqlPlans
}
