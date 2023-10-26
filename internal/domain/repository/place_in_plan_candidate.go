package repository

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceInPlanCandidateRepository interface {
	Save(ctx context.Context, planCandidateId string, place models.PlaceInPlanCandidate) error

	SavePlaces(ctx context.Context, planCandidateId string, places []models.PlaceInPlanCandidate) error

	FindByPlanCandidateId(ctx context.Context, planCandidateId string) (*[]models.PlaceInPlanCandidate, error)

	SaveGoogleImages(ctx context.Context, planCandidateId string, googlePlaceId string, images []models.Image) error

	SaveGoogleReviews(ctx context.Context, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error

	DeleteByPlanCandidateId(ctx context.Context, planCandidateId string) error
}
