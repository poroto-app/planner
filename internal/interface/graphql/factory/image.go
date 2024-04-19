package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/interface/graphql/model"
)

func ImageFromDomainModel(image *models.ImageSmallLarge) *graphql.Image {
	if image == nil {
		return nil
	}

	return &graphql.Image{
		Default: image.Default(),
		Small:   image.Small,
		Large:   image.Large,
		Google:  image.IsGooglePhotos,
	}
}
