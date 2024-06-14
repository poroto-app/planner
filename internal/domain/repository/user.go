package repository

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error

	Find(ctx context.Context, id string) (*models.User, error)

	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*models.User, error)

	UpdateProfile(ctx context.Context, userId string, name *string, photoUrl *string) error
}
