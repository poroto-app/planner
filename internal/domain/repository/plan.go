package repository

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

type SortedByCreatedAtQueryCursor string

type PlanRepository interface {
	Save(ctx context.Context, plan *models.Plan) error

	SortedByCreatedAt(ctx context.Context, queryCursor *SortedByCreatedAtQueryCursor, limit int) (*[]models.Plan, *SortedByCreatedAtQueryCursor, error)

	Find(ctx context.Context, planId string) (*models.Plan, error)

	FindByAuthorId(ctx context.Context, authorId string) (*[]models.Plan, error)

	// SortedByLocation location で指定した地点に近いプランを返す
	SortedByLocation(ctx context.Context, location models.GeoLocation, queryCursor *string, limit int) (*[]models.Plan, *string, error)

	// UpdatePlanAuthorUserByPlanCandidateSet プラン候補に紐づくプランの作者をユーザーに紐づける
	UpdatePlanAuthorUserByPlanCandidateSet(ctx context.Context, userId string, planCandidateSetIds []string) error

	FindCollage(ctx context.Context, planId string) (*models.PlanCollage, error)

	UpdateCollageImage(ctx context.Context, planId string, placeId string, placePhotoUrl string) error
}
