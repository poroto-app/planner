package user

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

type CheckUserAuthStateInput struct {
	UserId            string
	FirebaseAuthToken string
}

type CheckUserAuthStateOutput struct {
	IsAuthenticated bool
	User            *models.User
}

// CheckUserAuthState ユーザーの認証状態を確認する。
func (s Service) CheckUserAuthState(
	ctx context.Context,
	input CheckUserAuthStateInput,
) (*CheckUserAuthStateOutput, error) {
	user, err := s.FindByUserId(ctx, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("error while finding user by id: %v", err)
	}

	validUser, err := s.firebaseAuth.Verify(ctx, user.FirebaseUID, input.FirebaseAuthToken)
	if err != nil {
		return nil, fmt.Errorf("error while verifying firebase auth: %v", err)
	}

	return &CheckUserAuthStateOutput{
		IsAuthenticated: validUser,
		User:            user,
	}, nil
}
