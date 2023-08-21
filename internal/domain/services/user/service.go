package user

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/auth"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	userRepository repository.UserRepository
	firebaseAuth   *auth.FirebaseAuth
}

func NewService(ctx context.Context) (*Service, error) {
	userRepository, err := firestore.NewUserRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing user repository: %v", err)
	}

	firebaseAuth, err := auth.NewFirebaseAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firebase auth: %v", err)
	}

	return &Service{
		userRepository: userRepository,
		firebaseAuth:   firebaseAuth,
	}, nil
}
