package repository

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateRepository interface {
	// Create プラン候補を作成する
	// この時点ではプランは保存されない
	Create(cxt context.Context, planCandidateSetId string, expiresAt time.Time) error

	Find(ctx context.Context, planCandidateSetId string, now time.Time) (*models.PlanCandidateSet, error)

	FindPlan(ctx context.Context, planCandidateSetId string, planId string) (*models.Plan, error)

	FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]string, error)

	// AddPlan プラン候補にプランを追加する
	// 事前に models.PlanCandidateSet が保存されている必要がある
	AddPlan(ctx context.Context, planCandidateSetId string, plans ...models.Plan) error

	AddPlaceToPlan(ctx context.Context, planCandidateSetId string, planId string, previousPlaceId string, place models.Place) error

	RemovePlaceFromPlan(ctx context.Context, planCandidateSetId string, planId string, placeId string) error

	UpdatePlacesOrder(ctx context.Context, planId string, planCandidateSetId string, placeIdsOrdered []string) error

	UpdatePlanCandidateMetaData(ctx context.Context, planCandidateSetId string, meta models.PlanCandidateMetaData) error

	UpdateIsPlaceSearched(ctx context.Context, planCandidateSetId string, isPlaceSearched bool) error

	ReplacePlace(ctx context.Context, planCandidateSetId string, planId string, placeIdToBeReplaced string, placeToReplace models.Place) error

	DeleteAll(ctx context.Context, planCandidateSetIds []string) error

	// TODO: PlaceRepository に移動する
	UpdateLikeToPlaceInPlanCandidateSet(ctx context.Context, planCandidateSetId string, placeId string, like bool) error
}
