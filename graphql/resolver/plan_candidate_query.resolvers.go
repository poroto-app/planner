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
	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
)

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
		Plans:                         factory.PlansFromDomainModel(&planCandidate.Plans, planCandidate.MetaData.LocationStart),
		CreatedBasedOnCurrentLocation: planCandidate.MetaData.CreatedBasedOnCurrentLocation,
	}, nil
}

// MatchInterests is the resolver for the matchInterests field.
func (r *queryResolver) MatchInterests(ctx context.Context, input *model.MatchInterestsInput) (*model.InterestCandidate, error) {
	service, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan candidate service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	createPlanSessionId := uuid.New().String()

	categoriesSearched, err := service.CategoriesNearLocation(
		ctx,
		models.GeoLocation{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
		createPlanSessionId,
	)
	if err != nil {
		log.Printf("error while searching categories for session[%s]: %v", createPlanSessionId, err)
		return nil, fmt.Errorf("internal server error")
	}

	var categories = []*model.LocationCategory{}
	for _, categorySearched := range categoriesSearched {
		categories = append(categories, &model.LocationCategory{
			Name:            categorySearched.Name,
			DisplayName:     categorySearched.DisplayName,
			Photo:           categorySearched.Photo,
			DefaultPhotoURL: categorySearched.DefaultPhoto,
		})
	}

	return &model.InterestCandidate{
		Session:    createPlanSessionId,
		Categories: categories,
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

// PlacesToAddForPlanCandidate is the resolver for the placesToAddForPlanCandidate field.
func (r *queryResolver) PlacesToAddForPlanCandidate(ctx context.Context, input model.PlacesToAddForPlanCandidateInput) (*model.PlacesToAddForPlanCandidateOutput, error) {
	s, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan candidate service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	// TODO: 指定されたプランIDが不正だった場合の対処をする
	placesToAdd, err := s.FetchPlacesToAdd(ctx, input.PlanCandidateID, input.PlanID, 10)
	if err != nil {
		log.Println("error while fetching places to add: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	var places []*model.Place
	for _, place := range placesToAdd {
		places = append(places, factory.PlaceFromDomainModel(&place))
	}

	return &model.PlacesToAddForPlanCandidateOutput{
		Places: places,
	}, nil
}

// PlacesToReplaceForPlanCandidate is the resolver for the placesToReplaceForPlanCandidate field.
func (r *queryResolver) PlacesToReplaceForPlanCandidate(ctx context.Context, input model.PlacesToReplaceForPlanCandidateInput) (*model.PlacesToReplaceForPlanCandidateOutput, error) {
	panic(fmt.Errorf("not implemented: PlacesToReplaceForPlanCandidate - placesToReplaceForPlanCandidate"))
}
