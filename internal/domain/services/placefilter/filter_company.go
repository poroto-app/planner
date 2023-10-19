package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"strings"
)

// FilterCompany 会社の場所をフィルタリングする
func FilterCompany(placesToFiler []models.GooglePlace) []models.GooglePlace {
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

	return FilterPlaces(placesToFiler, func(place models.GooglePlace) bool {
		for _, tag := range companyTags {
			if strings.Contains(place.Name, tag) {
				return false
			}
		}
		return true
	})
}
