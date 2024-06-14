package user

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/apperrors"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

type UpdateUserProfileInput struct {
	AuthorizedUser  models.User
	UserId          string
	Name            string
	ProfileImageUrl string
}

// UpdateUserProfile ユーザーのプロフィールを更新する
// 認証ユーザーと更新対象が一致しない場合は apperrors.ErrUnauthorized を返す
// 更新対象が空白の場合は更新を行わない
func (s Service) UpdateUserProfile(ctx context.Context, input UpdateUserProfileInput) (*models.User, error) {
	if input.AuthorizedUser.Id != input.UserId {
		return nil, apperrors.ErrUnauthorized
	}

	// 空白だけの場合は更新を行わない
	if err := s.userRepository.UpdateProfile(
		ctx,
		input.UserId,
		utils.StrOmitWhitespace(input.Name),
		utils.StrOmitWhitespace(input.ProfileImageUrl),
	); err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	user, err := s.userRepository.Find(ctx, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}
