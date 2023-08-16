package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/graphql/factory"
	"poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/plan"
	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
	"poroto.app/poroto/planner/internal/domain/services/plangen"
)

// CreatePlanByLocation is the resolver for the createPlanByLocation field.
func (r *mutationResolver) CreatePlanByLocation(ctx context.Context, input model.CreatePlanByLocationInput) (*model.CreatePlanByLocationOutput, error) {
	planGenService, err := plangen.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan generator service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	planCandidateService, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan candidate service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	// TODO: 必須パラメータにする
	createBasedOnCurrentLocation := false
	if input.CreatedBasedOnCurrentLocation != nil {
		createBasedOnCurrentLocation = *input.CreatedBasedOnCurrentLocation
	}

	// TODO: sessionIDをリクエストに含めるようにする（二重で作成されないようにするため）
	session := uuid.New().String()
	plans, err := planGenService.CreatePlanByLocation(
		ctx,
		session,
		models.GeoLocation{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
		&input.CategoriesPreferred,
		&input.CategoriesDisliked,
		input.FreeTime,
		createBasedOnCurrentLocation,
	)
	if err != nil {
		log.Println(err)
	}

	if err := planCandidateService.SavePlanCandidate(ctx, session, *plans, *input.CreatedBasedOnCurrentLocation); err != nil {
		log.Println("error while caching plan candidate: ", err)
	}

	return &model.CreatePlanByLocationOutput{
		Session: session,
		Plans:   factory.PlansFromDomainModel(plans),
	}, nil
}

// CreatePlanByPlace is the resolver for the createPlanByPlace field.
func (r *mutationResolver) CreatePlanByPlace(ctx context.Context, input model.CreatePlanByPlaceInput) (*model.CreatePlanByPlaceOutput, error) {
	planGenService, err := plangen.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan generator service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	planCreated, err := planGenService.CreatePlanFromPlace(
		ctx,
		input.Session,
		input.PlaceID,
	)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlan := factory.PlanFromDomainModel(*planCreated)
	return &model.CreatePlanByPlaceOutput{
		Plan: &graphqlPlan,
	}, nil
}

// ChangePlacesOrderInPlanCandidate is the resolver for the changePlacesOrderInPlanCandidate field.
func (r *mutationResolver) ChangePlacesOrderInPlanCandidate(ctx context.Context, input model.ChangePlacesOrderInPlanCandidateInput) (*model.ChangePlacesOrderInPlanCandidateOutput, error) {
	service, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	var currentLocation *model.GeoLocation
	if input.CurrentLatitude != nil && input.CurrentLongitude != nil {
		currentLocation = &model.GeoLocation{
			Latitude:  *input.CurrentLatitude,
			Longitude: *input.CurrentLongitude,
		}
	}

	planUpdated, err := service.ChangePlacesOrderPlanCandidate(ctx, input.PlanID, input.Session, input.PlaceIds, currentLocation)
	if err != nil {
		return nil, fmt.Errorf("could not change places order")
	}

	graphqlPlan := factory.PlanFromDomainModel(*planUpdated)
	return &model.ChangePlacesOrderInPlanCandidateOutput{
		Plan: &graphqlPlan,
	}, nil
}

// SavePlanFromCandidate is the resolver for the savePlanFromCandidate field.
func (r *mutationResolver) SavePlanFromCandidate(ctx context.Context, input model.SavePlanFromCandidateInput) (*model.SavePlanFromCandidateOutput, error) {
	service, err := plan.NewService(ctx)
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
	planService, err := plan.NewService(ctx)
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

// Plans is the resolver for the plans field.
func (r *queryResolver) Plans(ctx context.Context, pageKey *string) ([]*model.Plan, error) {
	service, err := plan.NewService(ctx)
	if err != nil {
		log.Println("error while initializing places api: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	plans, err := service.FetchPlans(ctx, pageKey)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("could not fetch plans")
	}

	return factory.PlansFromDomainModel(plans), nil
}

// PlansByLocation is the resolver for the plansByLocation field.
func (r *queryResolver) PlansByLocation(ctx context.Context, input model.PlansByLocationInput) (*model.PlansByLocationOutput, error) {
	planService, err := plan.NewService(ctx)
	if err != nil {
		log.Printf("error while initializing plan service: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	plans, nextPageToken, err := planService.FetchPlansByLocation(
		ctx,
		models.GeoLocation{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
		input.Limit,
		input.PageKey,
	)
	if err != nil {
		log.Printf("error while fetching plans by location: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.PlansByLocationOutput{
		Plans:   factory.PlansFromDomainModel(plans),
		PageKey: nextPageToken,
	}, nil
}

// MatchInterests is the resolver for the matchInterests field.
func (r *queryResolver) MatchInterests(ctx context.Context, input *model.MatchInterestsInput) (*model.InterestCandidate, error) {
	service, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan candidate service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	categoriesSearched, err := service.CategoriesNearLocation(
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
	planService, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan candidate service: ", err)
		return nil, fmt.Errorf("internal server error")
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

// AvailablePlacesForPlan is the resolver for the availablePlacesForPlan field.
func (r *queryResolver) AvailablePlacesForPlan(ctx context.Context, input model.AvailablePlacesForPlanInput) (*model.AvailablePlacesForPlan, error) {
	s, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan candidate service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	availablePlaces, err := s.FetchCandidatePlaces(ctx, input.Session)
	if err != nil {
		log.Println("error while fetching candidate places: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlaces := make([]*model.Place, len(availablePlaces))
	for i, place := range availablePlaces {
		graphqlPlaces[i] = factory.PlaceFromDomainModel(place)
	}

	return &model.AvailablePlacesForPlan{
		Places: graphqlPlaces,
	}, nil
}
