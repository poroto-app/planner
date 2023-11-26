package repository

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlanCandidateRepository interface {
	Save(cxt context.Context, planCandidate *models.PlanCandidate) error

	Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error)

	FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]string, error)

	AddPlan(ctx context.Context, planCandidateId string, plans ...models.Plan) error

	AddPlaceToPlan(ctx context.Context, planCandidateId string, planId string, previousPlaceId string, place models.Place) error

	RemovePlaceFromPlan(ctx context.Context, planCandidateId string, planId string, placeId string) error

	UpdatePlacesOrder(ctx context.Context, planId string, planCandidate string, placeIdsOrdered []string) (*models.Plan, error)

	ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeToReplace models.Place) error

	DeleteAll(ctx context.Context, planCandidateIds []string) error
}
