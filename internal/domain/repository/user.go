package repository

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error

	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*models.User, error)
}
