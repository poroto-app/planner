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
	"poroto.app/poroto/planner/internal/domain/services/user"
	"poroto.app/poroto/planner/internal/domain/utils"
	gcontext "poroto.app/poroto/planner/internal/interface/graphql/context"
	"poroto.app/poroto/planner/internal/interface/graphql/factory"
	"poroto.app/poroto/planner/internal/interface/graphql/generated"
	"poroto.app/poroto/planner/internal/interface/graphql/model"
)

// Plans is the resolver for the plans field.
func (r *userResolver) Plans(ctx context.Context, obj *model.User) ([]*model.Plan, error) {
	// TODO: N+1問題に対応する（https://blog.giftee.dev/2023-11-10-gqlgen-dataloader/）
	r.Logger.Info("User#Plans", zap.String("userId", obj.ID))

	author, err := r.UserService.FindByUserId(ctx, obj.ID)
	if err != nil {
		r.Logger.Error("error while fetching user by id", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	if author == nil {
		return nil, nil
	}

	plans, err := r.PlanService.PlansByUser(ctx, obj.ID)
	if err != nil {
		r.Logger.Error("error while fetching plans by user", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	return factory.PlansFromDomainModel(plans, nil), nil
}

// LikedPlaces is the resolver for the likedPlaces field.
func (r *userResolver) LikedPlaces(ctx context.Context, obj *model.User) ([]*model.Place, error) {
	r.Logger.Info("User#LikedPlaces", zap.String("userId", obj.ID))

	// ログインユーザーでない場合は空配列を返す
	authUser := gcontext.GetAuthUser(ctx)
	if authUser == nil {
		r.Logger.Debug("auth user is nil")
		return []*model.Place{}, nil
	}

	// ログインユーザーとリクエストユーザーが異なる場合は空配列を返す
	if authUser.Id != obj.ID {
		r.Logger.Debug(
			"auth user and request user are different",
			zap.String("authUserId", authUser.Id),
			zap.String("requestUserId", obj.ID),
		)
		return []*model.Place{}, nil
	}

	places, err := r.UserService.FindLikePlaces(ctx, user.FindLikedPlacesInput{
		UserId:    obj.ID,
		CheckAuth: utils.ToPointer(false),
	})
	if err != nil {
		r.Logger.Error("error while fetching liked places", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlaces := array.Map(*places, func(place models.Place) *model.Place {
		return factory.PlaceFromDomainModel(&place)
	})

	return graphqlPlaces, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }