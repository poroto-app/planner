package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/infrastructure/graphql/model"
)

func UserFromDomainModel(user *models.User) *graphql.User {
	if user == nil {
		return nil
	}

	return &graphql.User{
		ID:       user.Id,
		Name:     user.Name,
		PhotoURL: user.PhotoUrl,
	}
}
