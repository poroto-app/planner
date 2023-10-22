package repository

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type GooglePlaceSearchResultRepository interface {
	// SaveImagesIfNotExist すでに画像が保存されていなかった場合のみ、保存する
	SaveImagesIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, images []models.Image) error

	// SaveReviewsIfNotExist すでにレビューが保存されていなかった場合のみ、保存する
	SaveReviewsIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error
}
