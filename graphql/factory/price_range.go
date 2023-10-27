package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PriceRangeFromDomainModel(priceRange *models.PriceRange) *graphql.PriceRange {
	if priceRange == nil {
		return nil
	}

	return &graphql.PriceRange{
		PriceRangeMin:    priceRange.Min,
		PriceRangeMax:    priceRange.Max,
		GooglePriceLevel: priceRange.GooglePriceLevel,
	}
}
