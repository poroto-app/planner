package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.29

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/graphql/factory"
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

	session := uuid.New().String()

	if err := service.CachePlanCandidate(session, *plans); err != nil {
		log.Println("error while caching plan candidate: ", err)
	}

	return &model.CreatePlanByLocationOutput{
		Session: session,
		Plans:   factory.PlansFromDomainModel(plans),
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
	planService, err := services.NewPlanService()
	if err != nil {
		log.Println("error while initializing places api: ", err)
		return nil, err
	}

	planCandidate, err := planService.FindPlanCandidate(input.Session)
	if err != nil {
		log.Println("error while finding plan candidate: ", err)
		return nil, err
	}

	if planCandidate == nil {
		return &model.CachedCreatedPlans{
			Plans: nil,
		}, nil
	}

	return &model.CachedCreatedPlans{
		Plans: factory.PlansFromDomainModel(&planCandidate.Plans),
	}, nil
}
