package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

// FindOrCreateFirebaseUser firebaseUIDに対応するユーザーを取得する。
// 存在しない場合は、新規に作成する。
func (s Service) FindOrCreateFirebaseUser(
	ctx context.Context,
	firebaseUID string,
	token string,
) (*models.User, error) {
	validUser, err := s.firebaseAuth.Verify(ctx, firebaseUID, token)
	if err != nil {
		return nil, fmt.Errorf("error while verifying firebase auth: %v", err)
	}

	if !validUser {
		return nil, fmt.Errorf("invalid user")
	}

	user, err := s.userRepository.FindByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("error while finding user by firebase uid: %v", err)
	}

	if user != nil {
		return user, nil
	}

	// ユーザーが存在しない場合は、新規に作成する
	firebaseUser, err := s.firebaseAuth.GetUser(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("error while getting firebase user: %v", err)
	}

	user = &models.User{
		Id:          uuid.New().String(),
		FirebaseUID: firebaseUser.UID,
		Name:        firebaseUser.DisplayName,
		Email:       utils.StrOmitEmpty(firebaseUser.Email),
		PhotoUrl:    utils.StrOmitEmpty(firebaseUser.PhotoURL),
	}
	if err := s.userRepository.Create(ctx, *user); err != nil {
		return nil, fmt.Errorf("error while creating user: %v", err)
	}

	return user, nil
}
