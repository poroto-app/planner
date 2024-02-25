package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"poroto.app/poroto/planner/internal/domain/services/plan"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/interface/graphql/factory"

	"poroto.app/poroto/planner/internal/interface/graphql/model"
)

// UploadPlacePhotoInPlan is the resolver for the uploadPlacePhotoInPlan field.
func (r *mutationResolver) UploadPlacePhotoInPlan(ctx context.Context, inputs []*model.UploadPlacePhotoInPlanInput) (*model.UploadPlacePhotoInPlanOutput, error) {
	panic(fmt.Errorf("not implemented: UploadPlacePhotoInPlan - uploadPlacePhotoInPlan"))
}

// LikeToPlaceInPlan is the resolver for the likeToPlaceInPlan field.
func (r *mutationResolver) LikeToPlaceInPlan(ctx context.Context, input model.LikeToPlaceInPlanInput) (*model.LikeToPlaceInPlanOutput, error) {
	logger, err := utils.NewLogger(utils.LoggerOption{Tag: "GraphQL"})
	if err != nil {
		return nil, fmt.Errorf("error while creating logger: %v", err)
	}

	planService, err := plan.NewService(ctx, r.DB)
	if err != nil {
		return nil, fmt.Errorf("error while creating plan service: %v", err)
	}

	logger.Info(
		"LikeToPlaceInPlan",
		zap.String("planId", input.PlanID),
		zap.String("placeId", input.PlaceID),
		zap.String("userId", input.UserID),
		zap.Bool("like", input.Like),
	)

	plan, err := planService.LikeToPlace(
		ctx,
		plan.LikeToPlaceInput{
			PlanId:            input.PlanID,
			PlaceId:           input.PlaceID,
			Like:              input.Like,
			UserId:            input.UserID,
			FirebaseAuthToken: input.FirebaseAuthToken,
		})
	if err != nil {
		logger.Error("error while liking to place in plan", zap.Error(err))
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	graphqlPlan, err := factory.PlanFromDomainModel(*plan, nil)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.LikeToPlaceInPlanOutput{
		Plan: graphqlPlan,
	}, nil
}
