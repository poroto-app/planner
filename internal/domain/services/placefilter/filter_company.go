package placefilter

import (
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"strings"
)

// FilterCompany 会社の場所をフィルタリングする
func FilterCompany(placesToFiler []places.Place) []places.Place {
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

	return FilterPlaces(placesToFiler, func(place places.Place) bool {
		for _, tag := range companyTags {
			if strings.Contains(place.Name, tag) {
				return false
			}
		}
		return true
	})
}
