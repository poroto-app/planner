package factory

import (
	"github.com/volatiletech/null/v8"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewUserEntityFromUser(user models.User) generated.User {
	return generated.User{
		ID:          user.Id,
		FirebaseUID: user.FirebaseUID,
		Name:        null.StringFrom(user.Name),
		PhotoURL:    null.StringFrom(utils.StrEmptyIfNil(user.PhotoUrl)),
		Email:       null.StringFrom(utils.StrEmptyIfNil(user.Email)),
	}
}

func NewUserFromUserEntity(userEntity generated.User) *models.User {
	return &models.User{
		Id:          userEntity.ID,
		FirebaseUID: userEntity.FirebaseUID,
		Name:        userEntity.Name.String,
		PhotoUrl:    utils.StrOmitEmpty(userEntity.PhotoURL.String),
		Email:       utils.StrOmitEmpty(userEntity.Email.String),
	}
}
