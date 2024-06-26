package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/domain/services/plan"
	"poroto.app/poroto/planner/internal/interface/graphql/factory"
	"poroto.app/poroto/planner/internal/interface/graphql/model"
)

// UploadPlacePhotoInPlan is the resolver for the uploadPlacePhotoInPlan field.
func (r *mutationResolver) UploadPlacePhotoInPlan(ctx context.Context, planID string, userID string, firebaseAuthToken string, inputs []*model.UploadPlacePhotoInPlanInput) (*model.UploadPlacePhotoInPlanOutput, error) {
	var uploadPlacePhotoInPlanInputs []place.UploadPlacePhotoInPlanInput
	for _, input := range inputs {
		uploadPlacePhotoInPlanInputs = append(uploadPlacePhotoInPlanInputs, place.UploadPlacePhotoInPlanInput{
			PlaceId:  input.PlaceID,
			PhotoUrl: input.PhotoURL,
			Width:    input.Width,
			Height:   input.Height,
		})
	}
	err := r.PlaceService.UploadPlacePhotoInPlan(ctx, userID, firebaseAuthToken, uploadPlacePhotoInPlanInputs)
	if err != nil {
		r.Logger.Error("error while uploading place photo in plan", zap.Error(err))
		return nil, fmt.Errorf("internal resolver error")
	}

	planDomainModel, err := r.PlanService.FetchPlan(ctx, planID)
	if err != nil {
		r.Logger.Error("error while fetching plan", zap.Error(err))
		return nil, fmt.Errorf("internal resolver error")
	}

	planGraphQLModel, err := factory.PlanFromDomainModel(*planDomainModel, nil)
	if err != nil {
		r.Logger.Error("error while converting plan domain model to graphql model", zap.Error(err))
		return nil, fmt.Errorf("internal resolver error")
	}
	return &model.UploadPlacePhotoInPlanOutput{
		Plan: planGraphQLModel,
	}, nil
}

// LikeToPlaceInPlan is the resolver for the likeToPlaceInPlan field.
func (r *mutationResolver) LikeToPlaceInPlan(ctx context.Context, input model.LikeToPlaceInPlanInput) (*model.LikeToPlaceInPlanOutput, error) {
	r.Logger.Info(
		"LikeToPlaceInPlan",
		zap.String("planId", input.PlanID),
		zap.String("placeId", input.PlaceID),
		zap.String("userId", input.UserID),
		zap.Bool("like", input.Like),
	)

	likeToPlaceResult, err := r.PlanService.LikeToPlace(
		ctx,
		plan.LikeToPlaceInput{
			PlanId:            input.PlanID,
			PlaceId:           input.PlaceID,
			Like:              input.Like,
			UserId:            input.UserID,
			FirebaseAuthToken: input.FirebaseAuthToken,
		})
	if err != nil {
		r.Logger.Error("error while liking to place in plan", zap.Error(err))
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	graphqlPlan, err := factory.PlanFromDomainModel(likeToPlaceResult.Plan, nil)
	if err != nil {
		r.Logger.Error("error while converting plan domain model to graphql model", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	return &model.LikeToPlaceInPlanOutput{
		Plan: graphqlPlan,
		LikedPlaceIds: array.Map(likeToPlaceResult.LikePlacesByUser, func(p models.Place) string {
			return p.Id
		}),
	}, nil
}

// UpdatePlanCollageImage is the resolver for the updatePlanCollageImage field.
func (r *mutationResolver) UpdatePlanCollageImage(ctx context.Context, input model.UpdatePlanCollageImageInput) (*model.UpdatePlanCollageImageOutput, error) {
	output, err := r.PlanService.UpdatePlanCollageImage(ctx, plan.UpdatePlanCollageImageInput{
		PlanId:            input.PlanID,
		PlaceId:           input.PlaceID,
		ImageUrl:          input.ImageURL,
		UserId:            input.UserID,
		FirebaseAuthToken: input.FirebaseAuthToken,
	})
	if err != nil {
		r.Logger.Error("error while updating plan collage image", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlan, err := factory.PlanFromDomainModel(output.Plan, nil)
	if err != nil {
		r.Logger.Error("error while converting plan domain model to graphql model", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	return &model.UpdatePlanCollageImageOutput{
		Plan: graphqlPlan,
	}, nil
}
