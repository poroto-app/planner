package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/interface/graphql/model"
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
