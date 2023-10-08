package repository

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"

	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type PlaceSearchResultRepository interface {
	Save(ctx context.Context, planCandidateId string, places []places.Place) error

	Find(ctx context.Context, planCandidateId string) ([]places.Place, error)

	// SaveImagesIfNotExist すでに画像が保存されていなかった場合のみ、保存する
	SaveImagesIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, images []models.Image) error

	DeleteAll(ctx context.Context, planCandidateIds []string) error
}
