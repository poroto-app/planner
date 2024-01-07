package entity

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestGooglePlacePhotoEntity_ToGooglePlacePhotoReference(t *testing.T) {
	cases := []struct {
		name                   string
		googlePlacePhotoEntity GooglePlacePhotoEntity
		expected               models.GooglePlacePhotoReference
	}{
		{
			name: "success",
			googlePlacePhotoEntity: GooglePlacePhotoEntity{
				PhotoReference:   "photo_reference",
				Width:            100,
				Height:           200,
				HTMLAttributions: []string{"html_attributions"},
				Small:            nil,
				Large:            nil,
			},
			expected: models.GooglePlacePhotoReference{
				PhotoReference:   "photo_reference",
				Width:            100,
				Height:           200,
				HTMLAttributions: []string{"html_attributions"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.googlePlacePhotoEntity.ToGooglePlacePhotoReference()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}
