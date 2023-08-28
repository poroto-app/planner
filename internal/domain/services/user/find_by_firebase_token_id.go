package user

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) FindByFirebaseIdToken(
	ctx context.Context,
	firebaseIdToken string,
) (*models.User, error) {
	firebaseUid, err := s.firebaseAuth.GetFirebaseUIDFromTokenId(ctx, firebaseIdToken)
	if err != nil {
		return nil, fmt.Errorf("error while getting firebase uid from token: %v", err)
	}

	user, err := s.userRepository.FindByFirebaseUID(ctx, *firebaseUid)
	if err != nil {
		return nil, fmt.Errorf("error while finding user by firebase uid: %v", err)
	}

	return user, nil
}
