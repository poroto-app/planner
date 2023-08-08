package mock

import (
	"context"

	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type PlaceSearchResultRepository struct {
	Data map[string][]places.Place
}

func NewPlaceSearchResultRepository(data map[string][]places.Place) PlaceSearchResultRepository {
	return PlaceSearchResultRepository{
		Data: data,
	}
}

func (p PlaceSearchResultRepository) Save(ctx context.Context, planCandidateId string, places []places.Place) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceSearchResultRepository) Find(ctx context.Context, planCandidateId string) ([]places.Place, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlaceSearchResultRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	for _, id := range planCandidateIds {
		delete(p.Data, id)
	}
	return nil
}
