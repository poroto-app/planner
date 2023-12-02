package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"strings"
)

// FilterCompany 会社の場所をフィルタリングする
func FilterCompany(placesToFiler []models.Place) []models.Place {
	companyTags := []string{
		"（株）",
		"（有）",
		"（合）",
		"(株)",
		"(有)",
		"(合)",
		"株式会社",
		"有限会社",
		"合同会社",
	}

	return FilterPlaces(placesToFiler, func(place models.Place) bool {
		for _, tag := range companyTags {
			if strings.Contains(place.Google.Name, tag) {
				return false
			}
		}
		return true
	})
}
