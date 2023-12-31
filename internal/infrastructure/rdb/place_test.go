package rdb

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"testing"
)

func TestPlaceRepository_SavePlacesFromGooglePlace(t *testing.T) {
	cases := []struct {
		name        string
		googlePlace models.GooglePlace
	}{
		{
			name: "save places from google place with nearby search result",
			googlePlace: models.GooglePlace{
				PlaceId:  "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
				Name:     "東京駅",
				Location: models.GeoLocation{Latitude: 35.6812362, Longitude: 139.7649361},
				Types:    []string{"train_station", "transit_station", "point_of_interest", "establishment"},
				PhotoReferences: []models.GooglePlacePhotoReference{
					{
						PhotoReference:   "AWU5eFjiROQJEeMpt7Hh2Pv_PIYOPIYO-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
						Width:            4032,
						Height:           3024,
						HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					},
				},
				PriceLevel:       2,
				Rating:           4.300000190734863,
				UserRatingsTotal: 100,
				Vicinity:         utils.StrPointer("日本、〒100-0005 東京都千代田区丸の内１丁目９−１"),
				Photos:           nil,
				PlaceDetail:      nil,
			},
		},
		{
			name: "save places from google place with place detail result",
			googlePlace: models.GooglePlace{
				PlaceId:  "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
				Name:     "東京駅",
				Location: models.GeoLocation{Latitude: 35.6812362, Longitude: 139.7649361},
				Types:    []string{"train_station", "transit_station", "point_of_interest", "establishment"},
				PhotoReferences: []models.GooglePlacePhotoReference{
					{
						PhotoReference:   "AWU5eFjiROQJEeMpt7Hh2Pv_PIYOPIYO-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
						Width:            4032,
						Height:           3024,
						HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					},
				},
				PriceLevel:       2,
				Rating:           4.300000190734863,
				UserRatingsTotal: 100,
				Vicinity:         utils.StrPointer("日本、〒100-0005 東京都千代田区丸の内１丁目９−１"),
				Photos: &[]models.GooglePlacePhoto{
					{
						PhotoReference:   "AWU5eFjiROQJEeMpt7Hh2Pv_PIYOPIYO-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
						Width:            4032,
						Height:           3024,
						HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
						Small:            utils.StrPointer("https://lh3.googleusercontent.com/places/HOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHOGEHGOEHOGEHOGEHOH=s1600-w1000-h1000"),
						Large:            utils.StrPointer("https://lh3.googleusercontent.com/places/FFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUGAFUG=s1600-w1000-h1000"),
					},
					{

						PhotoReference:   "AWU5eFgYAi-FUGAFUGA-lHUN-8Cbcl2xGP49EwZ5xzfo10jvcvuegwztrqV1iJmAjtG0XVs8Ph52lfav7mROP2Srh7h74CMNtXsQBKhIdFsjLp03zOcpfAWNkHqi4H54hyJ3VekpHvbiWOrayPbhnmWchlB5sLwcn17snJQ2uWA",
						Width:            4032,
						Height:           3024,
						HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>"},
						Small:            utils.StrPointer("https://lh3.googleusercontent.com/places/PPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYO=s1600-w1000-h1000"),
						Large:            utils.StrPointer("https://lh3.googleusercontent.com/places/PPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYOPIYO=s1600-w3000-h3000"),
					},
				},
				PlaceDetail: &models.GooglePlaceDetail{
					OpeningHours: &models.GooglePlaceOpeningHours{
						Periods: []models.GooglePlaceOpeningPeriod{
							{DayOfWeekOpen: "Monday", DayOfWeekClose: "Monday", OpeningTime: "1030", ClosingTime: "2130"},
							{DayOfWeekOpen: "Tuesday", DayOfWeekClose: "Tuesday", OpeningTime: "1030", ClosingTime: "2130"},
							{DayOfWeekOpen: "Wednesday", DayOfWeekClose: "Wednesday", OpeningTime: "1030", ClosingTime: "2130"},
							{DayOfWeekOpen: "Thursday", DayOfWeekClose: "Thursday", OpeningTime: "1030", ClosingTime: "2130"},
							{DayOfWeekOpen: "Friday", DayOfWeekClose: "Friday", OpeningTime: "1030", ClosingTime: "2130"},
							{DayOfWeekOpen: "Saturday", DayOfWeekClose: "Saturday", OpeningTime: "1030", ClosingTime: "2130"},
						},
					},
					PhotoReferences: []models.GooglePlacePhotoReference{
						{

							PhotoReference:   "AWU5eFgYAi-FUGAFUGA-lHUN-8Cbcl2xGP49EwZ5xzfo10jvcvuegwztrqV1iJmAjtG0XVs8Ph52lfav7mROP2Srh7h74CMNtXsQBKhIdFsjLp03zOcpfAWNkHqi4H54hyJ3VekpHvbiWOrayPbhnmWchlB5sLwcn17snJQ2uWA",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>"},
						},
					},
					Reviews: []models.GooglePlaceReview{
						{
							Rating:                4,
							Text:                  utils.StrPointer("とても大きな駅です。地下街も広く、お店もたくさんあります。駅員さんも多く、親切です。"),
							Time:                  1648126226,
							AuthorName:            "Alice Alicia",
							AuthorProfileImageUrl: utils.StrPointer("https://lh3.googleusercontent.com/a/ACg8ocKaPr9FWIiqs88c_Fugafugafugafugagfuagaufaugafufa=s128-c0x00000000-cc-rp-mo-ba5"),
							AuthorUrl:             utils.StrPointer("https://www.google.com/maps/contrib/117028493732372946396/reviews"),
						},
						{
							Rating:                5,
							Text:                  utils.StrPointer("近くに住んでいるので、よく利用しています。駅員さんも親切で、地下街も広く、お店もたくさんあります。"),
							Time:                  1618085426,
							AuthorName:            "Bob Bobson",
							AuthorProfileImageUrl: utils.StrPointer("https://lh3.googleusercontent.com/a-/ALV-HOGEhogehogehoge_wD8wQ5y5NPqCU7qZM9rnp00GHZYagec=s128-c0x00000000-cc-rp-mo-ba4"),
							AuthorUrl:             utils.StrPointer("https://www.google.com/maps/contrib/2849473937494373893093/reviews"),
						},
					},
				},
			},
		},
	}

	if testDB == nil {
		t.Fatalf("testDB is nil")
	}

	placeRepository, err := NewPlaceRepository(testDB)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer cleanup(context.Background(), testDB)

			_, err = placeRepository.SavePlacesFromGooglePlace(context.Background(), c.googlePlace)
			if err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			// GooglePlace が保存されているか確認
			isGooglePlaceSaved, err := entities.
				GooglePlaces(entities.GooglePlaceWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
				Exists(context.Background(), testDB)
			if err != nil {
				t.Fatalf("error while checking google place existence: %v", err)
			}
			if !isGooglePlaceSaved {
				t.Fatalf("google place is not saved")
			}

			// GooglePlaceType が保存されているか確認
			placeTypeCount, err := entities.
				GooglePlaceTypes(entities.GooglePlaceTypeWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
				Count(context.Background(), testDB)
			if err != nil {
				t.Fatalf("error while counting place types: %v", err)
			}

			if int(placeTypeCount) != len(c.googlePlace.Types) {
				t.Fatalf("place type expected: %d, actual: %d", len(c.googlePlace.Types), placeTypeCount)
			}

			// GooglePhotoReference が保存されているか確認
			for _, photoReference := range c.googlePlace.PhotoReferences {
				isPhotoReferenceSaved, err := entities.
					GooglePlacePhotoReferences(entities.GooglePlacePhotoReferenceWhere.PhotoReference.EQ(photoReference.PhotoReference)).
					Exists(context.Background(), testDB)
				if err != nil {
					t.Fatalf("error while checking photo reference existence: %v", err)
				}
				if !isPhotoReferenceSaved {
					t.Fatalf("photo is not saved")
				}
			}

			// HTMLAttributions が保存されているか確認
			for _, photoReference := range c.googlePlace.PhotoReferences {
				htmlAttributionCount, err := entities.
					GooglePlacePhotoAttributions(entities.GooglePlacePhotoAttributionWhere.PhotoReference.EQ(photoReference.PhotoReference)).
					Count(context.Background(), testDB)
				if err != nil {
					t.Fatalf("error while counting html attributions: %v", err)
				}

				if int(htmlAttributionCount) != len(photoReference.HTMLAttributions) {
					t.Fatalf("html attribution expected: %d, actual: %d", len(photoReference.HTMLAttributions), htmlAttributionCount)
				}
			}

			// Photo が保存されているか確認
			if c.googlePlace.Photos != nil {
				for _, photo := range *c.googlePlace.Photos {
					// 大・小の２つのバリエーションが保存されているか確認
					photoVariation := 0
					if photo.Small != nil {
						photoVariation++
					}
					if photo.Large != nil {
						photoVariation++
					}

					photoCount, err := entities.
						GooglePlacePhotos(entities.GooglePlacePhotoWhere.PhotoReference.EQ(photo.PhotoReference)).
						Count(context.Background(), testDB)
					if err != nil {
						t.Fatalf("error while counting google photos: %v", err)
					}
					if int(photoCount) != photoVariation {
						t.Fatalf("google photo expected: %d, actual: %d", photoVariation, photoCount)
					}
				}
			}

			if c.googlePlace.PlaceDetail != nil {
				// GooglePlaceOpeningPeriods が保存されているか確認
				if c.googlePlace.PlaceDetail.OpeningHours != nil {
					openingPeriodCount, err := entities.
						GooglePlaceOpeningPeriods(entities.GooglePlaceOpeningPeriodWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
						Count(context.Background(), testDB)
					if err != nil {
						t.Fatalf("error while counting opening periods: %v", err)
					}

					if int(openingPeriodCount) != len(c.googlePlace.PlaceDetail.OpeningHours.Periods) {
						t.Fatalf("opening period expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.OpeningHours.Periods), openingPeriodCount)
					}
				}

				// GooglePlaceReviews が保存されているか確認
				reviewCount, err := entities.
					GooglePlaceReviews(entities.GooglePlaceReviewWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
					Count(context.Background(), testDB)
				if err != nil {
					t.Fatalf("error while counting reviews: %v", err)
				}

				if int(reviewCount) != len(c.googlePlace.PlaceDetail.Reviews) {
					t.Fatalf("review expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.Reviews), reviewCount)
				}

				// GooglePhotoReference が保存されているか確認
				for _, photoReference := range c.googlePlace.PlaceDetail.PhotoReferences {
					isPhotoReferenceSaved, err := entities.
						GooglePlacePhotoReferences(entities.GooglePlacePhotoReferenceWhere.PhotoReference.EQ(photoReference.PhotoReference)).
						Exists(context.Background(), testDB)
					if err != nil {
						t.Fatalf("error while checking photo reference existence: %v", err)
					}
					if !isPhotoReferenceSaved {
						t.Fatalf("photo is not saved")
					}
				}

				// HTMLAttributions が保存されているか確認
				for _, photoReference := range c.googlePlace.PlaceDetail.PhotoReferences {
					htmlAttributionCount, err := entities.
						GooglePlacePhotoAttributions(entities.GooglePlacePhotoAttributionWhere.PhotoReference.EQ(photoReference.PhotoReference)).
						Count(context.Background(), testDB)
					if err != nil {
						t.Fatalf("error while counting html attributions: %v", err)
					}

					if int(htmlAttributionCount) != len(photoReference.HTMLAttributions) {
						t.Fatalf("html attribution expected: %d, actual: %d", len(photoReference.HTMLAttributions), htmlAttributionCount)
					}
				}
			}
		})
	}
}
