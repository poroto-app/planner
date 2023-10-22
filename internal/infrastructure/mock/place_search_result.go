package mock

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceSearchResultRepository struct {
	Data map[string][]models.GooglePlace
}

func NewPlaceSearchResultRepository(data map[string][]models.GooglePlace) PlaceSearchResultRepository {
	return PlaceSearchResultRepository{
		Data: data,
	}
}

func (p PlaceSearchResultRepository) Save(ctx context.Context, planCandidateId string, places []models.GooglePlace) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceSearchResultRepository) Find(ctx context.Context, planCandidateId string) ([]models.GooglePlace, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlaceSearchResultRepository) saveImagesIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, images []models.Image) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceSearchResultRepository) SaveReviewsIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceSearchResultRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	for _, id := range planCandidateIds {
		delete(p.Data, id)
	}
	return nil
}
