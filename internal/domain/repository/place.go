package repository

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceRepository interface {
	Find(ctx context.Context, placeId string) (*models.Place, error)

	// SavePlacesFromGooglePlaces はGooglePlaceからPlaceを作成し、保存する
	// すでに models.GooglePlace が保存されている場合は、それに紐づく models.Place を取得する
	SavePlacesFromGooglePlaces(ctx context.Context, googlePlaces ...models.GooglePlace) (*[]models.Place, error)

	FindByLocation(ctx context.Context, location models.GeoLocation, radius float64) ([]models.Place, error)

	// FindByGooglePlaceType は GooglePlaceType に紐づく Place を取得する
	FindByGooglePlaceType(ctx context.Context, googlePlaceType string, baseLocation models.GeoLocation, radius float64) (*[]models.Place, error)

	FindByGooglePlaceID(ctx context.Context, googlePlaceID string) (*models.Place, error)

	// FindLikePlacesByUserId はユーザーがいいねした Place を取得する
	FindLikePlacesByUserId(ctx context.Context, userId string) (*[]models.Place, error)

	// FindRecommendPlacesForCreatePlan は場所を指定してプランを作成するときに、おすすめの場所を取得する
	FindRecommendPlacesForCreatePlan(ctx context.Context) (*[]models.Place, error)

	SaveGooglePlacePhotos(ctx context.Context, googlePlaceId string, photos []models.GooglePlacePhoto) error

	SaveGooglePlaceDetail(ctx context.Context, googlePlaceId string, detail models.GooglePlaceDetail) error

	SavePlacePhotos(ctx context.Context, photos []models.PlacePhoto) error

	UpdateLikeByUserId(ctx context.Context, userId string, placeId string, like bool) error

	// UpdateLikeByPlanCandidateSetToUser PlanCandidateSet によりLikeされたものを、UserによるLikeに変更する
	UpdateLikeByPlanCandidateSetToUser(ctx context.Context, userId string, planCandidateSetIds []string) error
}
