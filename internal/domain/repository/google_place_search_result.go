package repository

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

type GooglePlaceSearchResultRepository interface {
	Save(ctx context.Context, planCandidateId string, places []models.GooglePlace) error

	Find(ctx context.Context, planCandidateId string) ([]models.GooglePlace, error)

	// SaveImagesIfNotExist すでに画像が保存されていなかった場合のみ、保存する
	SaveImagesIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, images []models.Image) error

	// SaveReviewsIfNotExist すでにレビューが保存されていなかった場合のみ、保存する
	SaveReviewsIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error

	// SavePriceLevel 上書き保存をする
	SavePriceLevel(ctx context.Context, planCandidateId string, googlePlaceId string, priceLevel *int) error

	DeleteAll(ctx context.Context, planCandidateIds []string) error
}
