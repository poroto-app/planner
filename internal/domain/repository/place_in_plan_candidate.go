package repository

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceInPlanCandidateRepository interface {
	SavePlaces(ctx context.Context, planCandidateId string, places []models.PlaceInPlanCandidate) error

	FindByGooglePlaceId(ctx context.Context, planCandidateId string, googlePlaceId string) (*models.PlaceInPlanCandidate, error)

	FindByPlanCandidateId(ctx context.Context, planCandidateId string) (*[]models.PlaceInPlanCandidate, error)

	SaveGooglePlacePhotos(ctx context.Context, planCandidateId string, googlePlaceId string, photos []models.GooglePlacePhoto) error

	SaveGooglePlaceDetail(ctx context.Context, planCandidateId string, googlePlaceId string, googlePlaceDetail models.GooglePlaceDetail) error

	DeleteByPlanCandidateId(ctx context.Context, planCandidateId string) error
}
