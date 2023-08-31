package user

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) FindByUserId(ctx context.Context, userId string) (*models.User, error) {
	user, err := s.userRepository.Find(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error while finding user by id: %v", err)
	}

	return user, nil
}
