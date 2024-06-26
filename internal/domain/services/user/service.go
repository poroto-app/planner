package user

import (
	"context"
	"database/sql"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/auth"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
)

type Service struct {
	userRepository  repository.UserRepository
	placeRepository repository.PlaceRepository
	planRepository  repository.PlanRepository
	firebaseAuth    *auth.FirebaseAuth
}

func NewService(ctx context.Context, db *sql.DB) (*Service, error) {
	userRepository, err := rdb.NewUserRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing user repository: %v", err)
	}

	placeRepository, err := rdb.NewPlaceRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place repository: %v", err)
	}

	planRepository, err := rdb.NewPlanRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan repository: %v", err)
	}

	firebaseAuth, err := auth.NewFirebaseAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firebase auth: %v", err)
	}

	return &Service{
		userRepository:  userRepository,
		placeRepository: placeRepository,
		planRepository:  planRepository,
		firebaseAuth:    firebaseAuth,
	}, nil
}
