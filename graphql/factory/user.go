package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
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
