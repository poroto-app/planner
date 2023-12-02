package entity

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"testing"
)

func TestGooglePlaceEntity_ToGooglePlace(t *testing.T) {
	cases := []struct {
		name              string
		googlePlaceEntity GooglePlaceEntity
		photoEntities     *[]GooglePlacePhotoEntity
		reviewEntities    *[]GooglePlaceReviewEntity
		expected          models.GooglePlace
	}{
		{
			name: "should return GooglePlace",
			googlePlaceEntity: GooglePlaceEntity{
				PlaceID: "place_id",
				Name:    "name",
				Types:   []string{"type1", "type2"},
				Location: GeoLocationEntity{
					Latitude:  1.0,
					Longitude: 2.0,
				},
				OpenNow:          true,
				Rating:           3.0,
				UserRatingsTotal: 10,
				PriceLevel:       2,
				OpeningHours: &GooglePlaceOpeningHoursEntity{
					OpeningHoursPeriods: []GooglePlaceOpeningPeriodEntity{
						{
							DayOfWeekOpen:  "Monday",
							DayOfWeekClose: "Monday",
							TimeOpen:       "10:00",
							TimeClose:      "20:00",
						},
					},
				},
			},
			photoEntities: &[]GooglePlacePhotoEntity{
				{
					PhotoReference: "photo_reference",
					Width:          100,
					Height:         200,
					Small:          utils.StrPointer("https://example.com/small"),
					Large:          utils.StrPointer("https://example.com/large"),
				},
			},
			reviewEntities: &[]GooglePlaceReviewEntity{
				{
					AuthorName: "author_name",
					AuthorUrl:  utils.StrPointer("https://example.com/author"),
					Language:   utils.StrPointer("ja"),
					Rating:     4.0,
				},
			},
			expected: models.GooglePlace{
				PlaceId:          "place_id",
				Name:             "name",
				Types:            []string{"type1", "type2"},
				Location:         models.GeoLocation{Latitude: 1.0, Longitude: 2.0},
				OpenNow:          true,
				Rating:           3.0,
				UserRatingsTotal: 10,
				PriceLevel:       2,
				Photos: &[]models.GooglePlacePhoto{
					{
						PhotoReference: "photo_reference",
						Width:          100,
						Height:         200,
						Small:          utils.StrPointer("https://example.com/small"),
						Large:          utils.StrPointer("https://example.com/large"),
					},
				},
				PlaceDetail: &models.GooglePlaceDetail{
					OpeningHours: &models.GooglePlaceOpeningHours{
						Periods: []models.GooglePlaceOpeningPeriod{
							{
								DayOfWeekOpen:  "Monday",
								DayOfWeekClose: "Monday",
								OpeningTime:    "10:00",
								ClosingTime:    "20:00",
							},
						},
					},
					Reviews: []models.GooglePlaceReview{
						{
							AuthorName: "author_name",
							AuthorUrl:  utils.StrPointer("https://example.com/author"),
							Language:   utils.StrPointer("ja"),
							Rating:     4.0,
						},
					},
					PhotoReferences: []models.GooglePlacePhotoReference{
						{
							PhotoReference: "photo_reference",
							Width:          100,
							Height:         200,
						},
					},
				},
			},
		},
		{
			name: "PlaceDetail is nil if photo, review and opening hours are nil",
			googlePlaceEntity: GooglePlaceEntity{
				PlaceID:      "place_id",
				OpeningHours: nil,
			},
			photoEntities:  nil,
			reviewEntities: nil,
			expected: models.GooglePlace{
				PlaceId:     "place_id",
				Photos:      nil,
				PlaceDetail: nil,
			},
		},
		{
			name: "PlaceDetail is nil if photo, review and opening hours are empty",
			googlePlaceEntity: GooglePlaceEntity{
				PlaceID:      "place_id",
				OpeningHours: nil,
			},
			photoEntities:  &[]GooglePlacePhotoEntity{},
			reviewEntities: &[]GooglePlaceReviewEntity{},
			expected: models.GooglePlace{
				PlaceId:     "place_id",
				Photos:      nil,
				PlaceDetail: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.googlePlaceEntity.ToGooglePlace(c.photoEntities, c.reviewEntities)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
