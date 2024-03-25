package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
	"poroto.app/poroto/planner/internal/interface/graphql/factory"
	"poroto.app/poroto/planner/internal/interface/graphql/model"
)

// PlanCandidate is the resolver for the planCandidate field.
func (r *queryResolver) PlanCandidate(ctx context.Context, input model.PlanCandidateInput) (*model.PlanCandidateOutput, error) {
	r.Logger.Info(
		"PlanCandidate",
		zap.String("planCandidateId", input.PlanCandidateID),
	)

	planCandidate, err := r.PlanCandidateService.FindPlanCandidate(ctx, plancandidate.FindPlanCandidateInput{
		PlanCandidateId:   input.PlanCandidateID,
		UserId:            input.UserID,
		FirebaseAuthToken: input.FirebaseAuthToken})
	if err != nil {
		r.Logger.Error("error while finding plan candidate", zap.Error(err))
		return nil, err
	}

	graphqlPlanCandidate := factory.PlanCandidateFromDomainModel(planCandidate)
	return &model.PlanCandidateOutput{
		PlanCandidate: graphqlPlanCandidate,
	}, nil
}

// NearbyPlaceCategories is the resolver for the nearbyPlaceCategories field.
func (r *queryResolver) NearbyPlaceCategories(ctx context.Context, input model.NearbyPlaceCategoriesInput) (*model.NearbyPlaceCategoryOutput, error) {
	createPlanSessionId := uuid.New().String()
	r.Logger.Info(
		"NearbyPlaceCategories",
		zap.String("planCandidateId", createPlanSessionId),
		zap.Float64("latitude", input.Latitude),
		zap.Float64("longitude", input.Longitude),
	)

	categoriesSearched, err := r.PlanCandidateService.CategoriesNearLocation(
		ctx,
		plancandidate.CategoryNearLocationParams{
			Location: models.GeoLocation{
				Latitude:  input.Latitude,
				Longitude: input.Longitude,
			},
			CreatePlanSessionId: createPlanSessionId,
		},
	)
	if err != nil {
		r.Logger.Error("error while finding nearby place categories", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	var categories []*model.NearbyLocationCategory
	for _, categorySearched := range categoriesSearched {
		var places []*model.Place
		for _, place := range categorySearched.Places {
			places = append(places, factory.PlaceFromDomainModel(&place))
		}

		categories = append(categories, &model.NearbyLocationCategory{
			ID:              categorySearched.Category.Name,
			DisplayName:     categorySearched.Category.DisplayName,
			DefaultPhotoURL: categorySearched.Category.DefaultPhoto,
			Places:          places,
		})
	}

	return &model.NearbyPlaceCategoryOutput{
		PlanCandidateID: createPlanSessionId,
		Categories:      categories,
	}, nil
}

// AvailablePlacesForPlan is the resolver for the availablePlacesForPlan field.
func (r *queryResolver) AvailablePlacesForPlan(ctx context.Context, input model.AvailablePlacesForPlanInput) (*model.AvailablePlacesForPlan, error) {
	availablePlaces, err := r.PlaceService.FetchCandidatePlaces(ctx, place.FetchCandidatePlacesInput{
		PlanCandidateId: input.Session,
	})
	if err != nil {
		r.Logger.Error("error while fetching available places for plan", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlaces := make([]*model.Place, len(*availablePlaces))
	for i, place := range *availablePlaces {
		graphqlPlaces[i] = factory.PlaceFromDomainModel(&place)
	}

	return &model.AvailablePlacesForPlan{
		Places: graphqlPlaces,
	}, nil
}

// PlacesToAddForPlanCandidate is the resolver for the placesToAddForPlanCandidate field.
func (r *queryResolver) PlacesToAddForPlanCandidate(ctx context.Context, input model.PlacesToAddForPlanCandidateInput) (*model.PlacesToAddForPlanCandidateOutput, error) {
	// TODO: 指定されたプランIDが不正だった場合の対処をする
	result, err := r.PlaceService.FetchPlacesToAdd(ctx, place.FetchPlacesToAddInput{
		PlanCandidateId: input.PlanCandidateID,
		PlanId:          input.PlanID,
		NLimit:          4,
	})
	if err != nil {
		r.Logger.Error("error while fetching places to add", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	var places []*model.Place
	for _, place := range result.PlacesRecommended {
		p := factory.PlaceFromDomainModel(&place)
		if p != nil {
			places = append(places, p)
		}
	}

	var placesGroupedByCategory []*model.CategoryGroupedPlaces
	for _, categoryGroupedPlaces := range result.PlacesGrouped {
		var places []*model.Place
		for _, place := range categoryGroupedPlaces.Places {
			p := factory.PlaceFromDomainModel(&place)
			if p != nil {
				places = append(places, p)
			}
		}
		placesGroupedByCategory = append(placesGroupedByCategory, &model.CategoryGroupedPlaces{
			Category: &model.PlaceCategory{
				ID:   categoryGroupedPlaces.Category.Name,
				Name: categoryGroupedPlaces.Category.DisplayName,
			},
			Places: places,
		})
	}

	return &model.PlacesToAddForPlanCandidateOutput{
		Places:                  places,
		PlacesGroupedByCategory: placesGroupedByCategory,
	}, nil
}

// PlacesToReplaceForPlanCandidate is the resolver for the placesToReplaceForPlanCandidate field.
func (r *queryResolver) PlacesToReplaceForPlanCandidate(ctx context.Context, input model.PlacesToReplaceForPlanCandidateInput) (*model.PlacesToReplaceForPlanCandidateOutput, error) {
	r.Logger.Info(
		"PlacesToReplaceForPlanCandidate",
		zap.String("planCandidateId", input.PlanCandidateID),
		zap.String("planId", input.PlanID),
		zap.String("placeId", input.PlaceID),
	)

	placesToReplace, err := r.PlaceService.FetchPlacesToReplace(ctx, input.PlanCandidateID, input.PlanID, input.PlaceID, 4)
	if err != nil {
		r.Logger.Error("error while fetching places to replace", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	var places []*model.Place
	for _, place := range placesToReplace {
		p := factory.PlaceFromDomainModel(&place)
		if p != nil {
			places = append(places, p)
		}
	}

	return &model.PlacesToReplaceForPlanCandidateOutput{
		Places: places,
	}, nil
}
