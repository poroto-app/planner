package repository

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceRepository interface {
	// SavePlacesFromGooglePlaces はGooglePlaceからPlaceを作成し、保存する
	// すでに models.GooglePlace が保存されている場合は、それに紐づく models.Place を取得する
	SavePlacesFromGooglePlaces(ctx context.Context, googlePlaces ...models.GooglePlace) (*[]models.Place, error)

	FindByLocation(ctx context.Context, location models.GeoLocation) ([]models.Place, error)

	FindByGooglePlaceID(ctx context.Context, googlePlaceID string) (*models.Place, error)

	// FindByPlanCandidateId は models.PlanCandidate に紐づく models.Place を取得する
	FindByPlanCandidateId(ctx context.Context, planCandidateId string) ([]models.Place, error)

	SaveGooglePlacePhotos(ctx context.Context, googlePlaceId string, photos []models.GooglePlacePhoto) error

	SaveGooglePlaceDetail(ctx context.Context, googlePlaceId string, detail models.GooglePlaceDetail) error
}
