package mock

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceInPlanCandidateRepository struct {
	Data map[string][]models.PlaceInPlanCandidate
}

func NewPlaceInPlanCandidateRepository(data map[string][]models.PlaceInPlanCandidate) *PlaceInPlanCandidateRepository {
	return &PlaceInPlanCandidateRepository{Data: data}
}

func (p PlaceInPlanCandidateRepository) Save(ctx context.Context, planCandidateId string, place models.PlaceInPlanCandidate) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceInPlanCandidateRepository) SavePlaces(ctx context.Context, planCandidateId string, places []models.PlaceInPlanCandidate) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceInPlanCandidateRepository) FindByPlanCandidateId(ctx context.Context, planCandidateId string) (*[]models.PlaceInPlanCandidate, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlaceInPlanCandidateRepository) SaveGooglePlacePhotos(ctx context.Context, planCandidateId string, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceInPlanCandidateRepository) SaveGoogleReviews(ctx context.Context, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceInPlanCandidateRepository) SaveGooglePlaceDetail(ctx context.Context, planCandidateId string, googlePlaceId string, googlePlaceDetail models.GooglePlaceDetail) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceInPlanCandidateRepository) DeleteByPlanCandidateId(ctx context.Context, planCandidateId string) error {
	delete(p.Data, planCandidateId)
	return nil
}
