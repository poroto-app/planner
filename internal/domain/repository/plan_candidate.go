package repository

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateRepository interface {
	// Create プラン候補を作成する
	// この時点ではプランは保存されない
	Create(cxt context.Context, planCandidateId string, expiresAt time.Time) error

	Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error)

	FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]string, error)

	// AddSearchedPlacesForPlanCandidate は models.PlanCandidate を作成するために検索した場所を保存する
	AddSearchedPlacesForPlanCandidate(ctx context.Context, planCandidateId string, placeIds []string) error

	// AddPlan プラン候補にプランを追加する
	// 事前に models.PlanCandidate が保存されている必要がある
	AddPlan(ctx context.Context, planCandidateId string, plans ...models.Plan) error

	AddPlaceToPlan(ctx context.Context, planCandidateId string, planId string, previousPlaceId string, place models.Place) error

	RemovePlaceFromPlan(ctx context.Context, planCandidateId string, planId string, placeId string) error

	UpdatePlacesOrder(ctx context.Context, planId string, planCandidate string, placeIdsOrdered []string) (*models.Plan, error)

	UpdatePlanCandidateMetaData(ctx context.Context, planCandidateId string, meta models.PlanCandidateMetaData) error

	ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeToReplace models.Place) error

	DeleteAll(ctx context.Context, planCandidateIds []string) error

	UpdateLikeToPlaceInPlanCandidate(ctx context.Context, planCandidateId string, planId string, placeId string, like bool) error
}
