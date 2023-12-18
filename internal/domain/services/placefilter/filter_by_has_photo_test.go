package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestFilterByHasPhoto(t *testing.T) {
	cases := []struct {
		name     string
		places   []models.Place
		expected []models.Place
	}{
		{
			name: "should return places which has photo",
			places: []models.Place{
				{
					Id: "has_photo",
					Google: models.GooglePlace{
						Photos: &[]models.GooglePlacePhoto{{PhotoReference: "photo_reference"}},
					},
				},
				{
					Id: "empty_photo",
					Google: models.GooglePlace{
						Photos: &[]models.GooglePlacePhoto{},
					},
				},
			},
			expected: []models.Place{
				{
					Id: "has_photo",
					Google: models.GooglePlace{
						Photos: &[]models.GooglePlacePhoto{{PhotoReference: "photo_reference"}},
					},
				},
			},
		},
		{
			name: "should return places which has photo reference",
			places: []models.Place{
				{
					Id: "has_photo_reference",
					Google: models.GooglePlace{
						PhotoReferences: []models.GooglePlacePhotoReference{{PhotoReference: "photo_reference"}},
					},
				},
				{
					Id: "empty_photo_reference",
					Google: models.GooglePlace{
						PhotoReferences: []models.GooglePlacePhotoReference{},
					},
				},
			},
			expected: []models.Place{
				{
					Id: "has_photo_reference",
					Google: models.GooglePlace{
						PhotoReferences: []models.GooglePlacePhotoReference{{PhotoReference: "photo_reference"}},
					},
				},
			},
		},
		{
			name: "should return places which has photo reference in place detail",
			places: []models.Place{
				{
					Id: "has_photo_reference_in_place_detail",
					Google: models.GooglePlace{
						PlaceDetail: &models.GooglePlaceDetail{
							PhotoReferences: []models.GooglePlacePhotoReference{{PhotoReference: "photo_reference"}},
						},
					},
				},
				{
					Id: "empty_photo_reference_in_place_detail",
					Google: models.GooglePlace{
						PlaceDetail: &models.GooglePlaceDetail{
							PhotoReferences: []models.GooglePlacePhotoReference{},
						},
					},
				},
			},
			expected: []models.Place{
				{
					Id: "has_photo_reference_in_place_detail",
					Google: models.GooglePlace{
						PlaceDetail: &models.GooglePlaceDetail{
							PhotoReferences: []models.GooglePlacePhotoReference{{PhotoReference: "photo_reference"}},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FilterByHasPhoto(c.places)
			if len(actual) != len(c.expected) {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
			for i, a := range actual {
				if a.Id != c.expected[i].Id {
					t.Errorf("expected: %v, actual: %v", c.expected, actual)
				}
			}
		})
	}
}
