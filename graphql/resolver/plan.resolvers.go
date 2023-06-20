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
	"poroto.app/poroto/planner/internal/domain/services/plan"
)

// CreatePlanByLocation is the resolver for the createPlanByLocation field.
func (r *mutationResolver) CreatePlanByLocation(ctx context.Context, input model.CreatePlanByLocationInput) (*model.CreatePlanByLocationOutput, error) {
	// TODO: エラー時は処理を停止させる
	service, err := plan.NewPlanService(ctx)
	if err != nil {
		log.Println(err)
	}

	plans, err := service.CreatePlanByLocation(
		ctx,
		models.GeoLocation{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
		&input.Categories,
		input.FreeTime)
	if err != nil {
		log.Println(err)
	}

	session := uuid.New().String()

	if err := service.CachePlanCandidate(ctx, session, *plans, *input.CreatedBasedOnCurrentLocation); err != nil {
		log.Println("error while caching plan candidate: ", err)
	}

	return &model.CreatePlanByLocationOutput{
		Session: session,
		Plans:   factory.PlansFromDomainModel(plans),
	}, nil
}

// ChangePlacesOrderInPlan is the resolver for the ChangePlacesOrderInPlan field.
func (r *mutationResolver) ChangePlacesOrderInPlan(ctx context.Context, input model.ChangePlacesOrderInPlanInput) (*model.ChangePlacesOrderInPlanOutput, error) {
	panic(fmt.Errorf("not implemented: ChangePlacesOrderInPlan - ChangePlacesOrderInPlan"))
}

// SavePlanFromCandidate is the resolver for the savePlanFromCandidate field.
func (r *mutationResolver) SavePlanFromCandidate(ctx context.Context, input model.SavePlanFromCandidateInput) (*model.SavePlanFromCandidateOutput, error) {
	service, err := plan.NewPlanService(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	planSaved, err := service.SavePlanFromPlanCandidate(ctx, input.Session, input.PlanID)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("could not save plan")
	}

	graphqlPlan := factory.PlanFromDomainModel(*planSaved)
	return &model.SavePlanFromCandidateOutput{
		Plan: &graphqlPlan,
	}, nil
}

// Plan is the resolver for the plan field.
func (r *queryResolver) Plan(ctx context.Context, id string) (*model.Plan, error) {
	planService, err := plan.NewPlanService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	p, err := planService.FetchPlan(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan: %v", err)
	}

	if p == nil {
		return nil, nil
	}

	graphqlPlan := factory.PlanFromDomainModel(*p)
	return &graphqlPlan, nil
}

// MatchInterests is the resolver for the matchInterests field.
func (r *queryResolver) MatchInterests(ctx context.Context, input *model.MatchInterestsInput) (*model.InterestCandidate, error) {
	planService, err := plan.NewPlanService(ctx)
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
	planService, err := plan.NewPlanService(ctx)
	if err != nil {
		log.Println("error while initializing places api: ", err)
		return nil, err
	}

	planCandidate, err := planService.FindPlanCandidate(ctx, input.Session)
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
		Plans:                         factory.PlansFromDomainModel(&planCandidate.Plans),
		CreatedBasedOnCurrentLocation: planCandidate.CreatedBasedOnCurrentLocation,
	}, nil
}
