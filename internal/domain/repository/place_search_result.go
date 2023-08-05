package repository

import (
	"context"

	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type PlaceSearchResultRepository interface {
	Save(ctx context.Context, planCandidateId string, places []places.Place) error
	Find(ctx context.Context, planCandidateId string) ([]places.Place, error)
}
