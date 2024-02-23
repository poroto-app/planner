package rdb

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
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
						PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
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
						PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
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
						PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
						Width:            4032,
						Height:           3024,
						HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
						Small: &models.Image{
							Width:  400,
							Height: 400,
							URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
						},
						Large: &models.Image{
							Width:  4032,
							Height: 3024,
							URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
						},
					},
					{

						PhotoReference:   "photo-2-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
						Width:            1920,
						Height:           1080,
						HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>"},
						Small: &models.Image{
							Width:  400,
							Height: 400,
							URL:    "https://lh3.googleusercontent.com/places/photo-2=s1600-w400-h400",
						},
						Large: &models.Image{
							Width:  1920,
							Height: 1080,
							URL:    "https://lh3.googleusercontent.com/places/photo-2=s1600-w1920-h1080",
						},
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
						// Place Detailで取得される値は一部、Nearby Searchで取得される値と重複する
						{
							PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
						},
						{
							PhotoReference:   "photo-2-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            1920,
							Height:           1080,
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
		{
			name: "save places from google place with duplicated values",
			googlePlace: models.GooglePlace{
				PlaceId:  "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
				Name:     "東京駅",
				Location: models.GeoLocation{Latitude: 35.6812362, Longitude: 139.7649361},
				Types:    []string{"train_station", "transit_station", "point_of_interest", "establishment"},
				PhotoReferences: []models.GooglePlacePhotoReference{
					{
						PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
						Width:            4032,
						Height:           3024,
						HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					},
					// 重複した値
					{
						PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
						Width:            4032,
						Height:           3024,
						HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					},
				},
			},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			defer func(ctx context.Context, db *sql.DB) {
				err := cleanup(ctx, db)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			}(testContext, testDB)

			firstSaveResult, err := placeRepository.SavePlacesFromGooglePlaces(testContext, c.googlePlace)
			if err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			// GooglePlace が保存されているか確認
			isGooglePlaceSaved, err := generated.
				GooglePlaces(generated.GooglePlaceWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
				Exists(testContext, testDB)
			if err != nil {
				t.Fatalf("error while checking google place existence: %v", err)
			}
			if !isGooglePlaceSaved {
				t.Fatalf("google place is not saved")
			}

			// GooglePlaceType が保存されているか確認
			placeTypeCount, err := generated.
				GooglePlaceTypes(generated.GooglePlaceTypeWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
				Count(testContext, testDB)
			if err != nil {
				t.Fatalf("error while counting place types: %v", err)
			}

			if int(placeTypeCount) != len(c.googlePlace.Types) {
				t.Fatalf("place type expected: %d, actual: %d", len(c.googlePlace.Types), placeTypeCount)
			}

			// GooglePhotoReference が保存されているか確認
			for _, photoReference := range c.googlePlace.PhotoReferences {
				isPhotoReferenceSaved, err := generated.
					GooglePlacePhotoReferences(generated.GooglePlacePhotoReferenceWhere.PhotoReference.EQ(photoReference.PhotoReference)).
					Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("error while checking photo reference existence: %v", err)
				}
				if !isPhotoReferenceSaved {
					t.Fatalf("photo is not saved")
				}
			}

			// HTMLAttributions が保存されているか確認
			for _, photoReference := range c.googlePlace.PhotoReferences {
				htmlAttributionCount, err := generated.
					GooglePlacePhotoAttributions(generated.GooglePlacePhotoAttributionWhere.PhotoReference.EQ(photoReference.PhotoReference)).
					Count(testContext, testDB)
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

					photoCount, err := generated.
						GooglePlacePhotos(generated.GooglePlacePhotoWhere.PhotoReference.EQ(photo.PhotoReference)).
						Count(testContext, testDB)
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
					openingPeriodCount, err := generated.
						GooglePlaceOpeningPeriods(generated.GooglePlaceOpeningPeriodWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
						Count(testContext, testDB)
					if err != nil {
						t.Fatalf("error while counting opening periods: %v", err)
					}

					if int(openingPeriodCount) != len(c.googlePlace.PlaceDetail.OpeningHours.Periods) {
						t.Fatalf("opening period expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.OpeningHours.Periods), openingPeriodCount)
					}
				}

				// GooglePlaceReviews が保存されているか確認
				reviewCount, err := generated.
					GooglePlaceReviews(generated.GooglePlaceReviewWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
					Count(testContext, testDB)
				if err != nil {
					t.Fatalf("error while counting reviews: %v", err)
				}

				if int(reviewCount) != len(c.googlePlace.PlaceDetail.Reviews) {
					t.Fatalf("review expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.Reviews), reviewCount)
				}

				// GooglePhotoReference が保存されているか確認
				for _, photoReference := range c.googlePlace.PlaceDetail.PhotoReferences {
					isPhotoReferenceSaved, err := generated.
						GooglePlacePhotoReferences(generated.GooglePlacePhotoReferenceWhere.PhotoReference.EQ(photoReference.PhotoReference)).
						Exists(testContext, testDB)
					if err != nil {
						t.Fatalf("error while checking photo reference existence: %v", err)
					}
					if !isPhotoReferenceSaved {
						t.Fatalf("photo is not saved")
					}
				}

				// HTMLAttributions が保存されているか確認
				for _, photoReference := range c.googlePlace.PlaceDetail.PhotoReferences {
					htmlAttributionCount, err := generated.
						GooglePlacePhotoAttributions(generated.GooglePlacePhotoAttributionWhere.PhotoReference.EQ(photoReference.PhotoReference)).
						Count(testContext, testDB)
					if err != nil {
						t.Fatalf("error while counting html attributions: %v", err)
					}

					if int(htmlAttributionCount) != len(photoReference.HTMLAttributions) {
						t.Fatalf("html attribution expected: %d, actual: %d", len(photoReference.HTMLAttributions), htmlAttributionCount)
					}
				}
			}

			// 一度保存したあとは、すでに保存されたものが取得される
			secondSaveResult, err := placeRepository.SavePlacesFromGooglePlaces(testContext, c.googlePlace)
			if err != nil {
				t.Fatalf("error while saving places second time: %v", err)
			}

			actualFirstSave := (*firstSaveResult)[0]
			actualSecondSave := (*secondSaveResult)[0]

			if actualFirstSave.Id != actualSecondSave.Id {
				t.Fatalf("place id expected: %s, actual: %s", actualFirstSave.Id, actualSecondSave.Id)
			}

			if len(actualFirstSave.Google.Types) != len(actualSecondSave.Google.Types) {
				t.Fatalf("place type expected: %d, actual: %d", len(actualFirstSave.Google.Types), len(actualSecondSave.Google.Types))
			}

			if len(actualFirstSave.Google.PhotoReferences) != len(actualSecondSave.Google.PhotoReferences) {
				t.Fatalf("photo reference expected: %d, actual: %d", len(actualFirstSave.Google.PhotoReferences), len(actualSecondSave.Google.PhotoReferences))
			}

			if c.googlePlace.Photos != nil {
				if len(*c.googlePlace.Photos) != len(*actualSecondSave.Google.Photos) {
					t.Fatalf("photo expected: %d, actual: %d", len(*c.googlePlace.Photos), len(*actualSecondSave.Google.Photos))
				}
			}

			if c.googlePlace.PlaceDetail != nil {
				if len(c.googlePlace.PlaceDetail.Reviews) != len(actualSecondSave.Google.PlaceDetail.Reviews) {
					t.Fatalf("review expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.Reviews), len(actualSecondSave.Google.PlaceDetail.Reviews))
				}

				if c.googlePlace.PlaceDetail.OpeningHours != nil {
					if len(c.googlePlace.PlaceDetail.OpeningHours.Periods) != len(actualSecondSave.Google.PlaceDetail.OpeningHours.Periods) {
						t.Fatalf("opening period expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.OpeningHours.Periods), len(actualSecondSave.Google.PlaceDetail.OpeningHours.Periods))
					}
				}
			}
		})
	}
}

func TestPlaceRepository_SavePlacesFromGooglePlace_DuplicatedValue(t *testing.T) {
	cases := []struct {
		name        string
		savedPlace  models.Place
		googlePlace models.GooglePlace
		expected    models.Place
	}{
		{
			name: "return saved google place if google place is already saved",
			savedPlace: models.Place{
				Id:       "aafd9600-c57d-494a-8f66-f4952f0fd475",
				Name:     "東京駅",
				Location: models.GeoLocation{Latitude: 35.6812362, Longitude: 139.7649361},
				Google: models.GooglePlace{
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
							Small: &models.Image{
								Width:  400,
								Height: 400,
								URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
							},
							Large: &models.Image{
								Width:  4032,
								Height: 3024,
								URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
							},
						},
						{

							PhotoReference:   "AWU5eFgYAi-FUGAFUGA-lHUN-8Cbcl2xGP49EwZ5xzfo10jvcvuegwztrqV1iJmAjtG0XVs8Ph52lfav7mROP2Srh7h74CMNtXsQBKhIdFsjLp03zOcpfAWNkHqi4H54hyJ3VekpHvbiWOrayPbhnmWchlB5sLwcn17snJQ2uWA",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>"},
							Small: &models.Image{
								Width:  400,
								Height: 400,
								URL:    "https://lh3.googleusercontent.com/places/photo-2=s1600-w400-h400",
							},
							Large: &models.Image{
								Width:  4032,
								Height: 3024,
								URL:    "https://lh3.googleusercontent.com/places/photo-2=s1600-w4032-h3024",
							},
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
			googlePlace: models.GooglePlace{
				PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
			},
			expected: models.Place{
				Id:       "aafd9600-c57d-494a-8f66-f4952f0fd475",
				Name:     "東京駅",
				Location: models.GeoLocation{Latitude: 35.6812362, Longitude: 139.7649361},
				Google: models.GooglePlace{
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
							Small: &models.Image{
								Width:  400,
								Height: 400,
								URL:    "https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w400-h400",
							},
							Large: &models.Image{
								Width:  4032,
								Height: 3024,
								URL:    "https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w4032-h3024",
							},
						},
						{

							PhotoReference:   "AWU5eFgYAi-FUGAFUGA-lHUN-8Cbcl2xGP49EwZ5xzfo10jvcvuegwztrqV1iJmAjtG0XVs8Ph52lfav7mROP2Srh7h74CMNtXsQBKhIdFsjLp03zOcpfAWNkHqi4H54hyJ3VekpHvbiWOrayPbhnmWchlB5sLwcn17snJQ2uWA",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>"},
							Small: &models.Image{
								Width:  400,
								Height: 400,
								URL:    "https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w400-h400",
							},
							Large: &models.Image{
								Width:  4032,
								Height: 3024,
								URL:    "https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w4032-h3024",
							},
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
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			defer func(ctx context.Context, db *sql.DB) {
				err := cleanup(ctx, db)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			}(testContext, testDB)

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, []models.Place{c.savedPlace}); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			result, err := placeRepository.SavePlacesFromGooglePlaces(testContext, c.googlePlace)
			if err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			actual := (*result)[0]

			// すでに保存されている Google Place が取得される
			if c.expected.Id != actual.Id {
				t.Fatalf("place id expected: %s, actual: %s", c.expected.Id, actual.Id)
			}
		})
	}
}

func TestPlaceRepository_FindByCategory(t *testing.T) {
	cases := []struct {
		name            string
		savedPlaces     []models.Place
		googlePlaceType string
		baseLocation    models.GeoLocation
		radius          float64
		expectedPlaces  []models.Place
	}{
		{
			name: "valid",
			savedPlaces: []models.Place{
				{
					Id: "kinokuniya-shoten",
					Google: models.GooglePlace{
						PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
						Location: models.GeoLocation{
							Latitude:  35.692247367825,
							Longitude: 139.703036771,
						},
						Types: []string{"book_store", "point_of_interest", "store", "establishment"},
					},
				},
			},
			googlePlaceType: "book_store",
			radius:          5000,
			baseLocation: models.GeoLocation{
				// 新宿駅
				Latitude:  35.6896,
				Longitude: 139.7005,
			},
			expectedPlaces: []models.Place{
				{
					Id: "kinokuniya-shoten",
					Location: models.GeoLocation{
						Latitude:  35.692247367825,
						Longitude: 139.703036771,
					},
					Google: models.GooglePlace{
						PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
						Location: models.GeoLocation{
							Latitude:  35.692247367825,
							Longitude: 139.703036771,
						},
						Types: []string{"book_store", "point_of_interest", "store", "establishment"},
					},
				},
			},
		},
		{
			name: "filter by googlePlaceType",
			savedPlaces: []models.Place{
				{
					Id: "kinokuniya-shoten",
					Google: models.GooglePlace{
						PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
						Location: models.GeoLocation{
							Latitude:  35.692247367825,
							Longitude: 139.703036771,
						},
						Types: []string{"book_store", "point_of_interest", "store", "establishment"},
					},
				},
			},
			googlePlaceType: "cafe",
			radius:          5000,
			baseLocation: models.GeoLocation{
				// 新宿駅
				Latitude:  35.6896,
				Longitude: 139.7005,
			},
			expectedPlaces: []models.Place{},
		},
		{
			name: "filter by distance",
			savedPlaces: []models.Place{
				{
					Id: "kinokuniya-shoten",
					Google: models.GooglePlace{
						PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
						Location: models.GeoLocation{
							Latitude:  35.692247367825,
							Longitude: 139.703036771,
						},
						Types: []string{"book_store", "point_of_interest", "store", "establishment"},
					},
				},
			},
			googlePlaceType: "book_store",
			radius:          1000,
			baseLocation: models.GeoLocation{
				// 代々木上原駅
				Latitude:  35.669017114155,
				Longitude: 139.67981467654,
			},
			expectedPlaces: []models.Place{},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			})

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			actualPlaces, err := placeRepository.FindByGooglePlaceType(testContext, c.googlePlaceType, c.baseLocation, c.radius)
			if err != nil {
				t.Fatalf("error while finding places: %v", err)
			}

			if diff := cmp.Diff(c.expectedPlaces, *actualPlaces); diff != "" {
				t.Fatalf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlaceRepository_FindByGooglePlaceID(t *testing.T) {
	cases := []struct {
		name          string
		savedPlaces   []models.Place
		googlePlaceId string
		expectedPlace *models.Place
	}{
		{
			name: "find place by google place id",
			savedPlaces: []models.Place{
				{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
				{Id: "place_id_2", Google: models.GooglePlace{PlaceId: "google_place_id_2"}},
			},
			googlePlaceId: "google_place_id_1",
			expectedPlace: &models.Place{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			})

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			actualPlace, err := placeRepository.FindByGooglePlaceID(testContext, c.googlePlaceId)
			if err != nil {
				t.Fatalf("error while finding place: %v", err)
			}

			if diff := cmp.Diff(c.expectedPlace, actualPlace); diff != "" {
				t.Fatalf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlaceRepository_FindByGooglePlaceID_WithLikeCount(t *testing.T) {
	cases := []struct {
		name                            string
		savedPlaces                     []models.Place
		savedPlanCandidateSets          []models.PlanCandidate
		savedPlanCandidateSetLikePlaces []generated.PlanCandidateSetLikePlace
		googlePlaceId                   string
		expectedPlace                   *models.Place
	}{
		{
			name: "find place by google place id with like count",
			savedPlaces: []models.Place{
				{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
			},
			savedPlanCandidateSets: []models.PlanCandidate{
				{Id: "plan_candidate_set_id_1", ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)},
			},
			savedPlanCandidateSetLikePlaces: []generated.PlanCandidateSetLikePlace{
				{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: "plan_candidate_set_id_1",
					PlaceID:            "place_id_1",
				},
			},
			googlePlaceId: "google_place_id_1",
			expectedPlace: &models.Place{
				Id:        "place_id_1",
				Google:    models.GooglePlace{PlaceId: "google_place_id_1"},
				LikeCount: 1,
			},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			})

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			// 事前にPlanCandidateSetを保存しておく
			for _, planCandidateSet := range c.savedPlanCandidateSets {
				if err := savePlanCandidate(testContext, testDB, planCandidateSet); err != nil {
					t.Fatalf("error while saving plan candidate set: %v", err)
				}
			}

			// 事前にPlanCandidateSetLikePlaceを保存しておく
			for _, planCandidateSetLikePlace := range c.savedPlanCandidateSetLikePlaces {
				if err := planCandidateSetLikePlace.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("error while saving plan candidate set like place: %v", err)
				}
			}

			actualPlace, err := placeRepository.FindByGooglePlaceID(testContext, c.googlePlaceId)
			if err != nil {
				t.Fatalf("error while finding place: %v", err)
			}

			if diff := cmp.Diff(c.expectedPlace, actualPlace); diff != "" {
				t.Fatalf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlaceRepository_FindByPlanCandidateId(t *testing.T) {
	cases := []struct {
		name                               string
		planCandidateId                    string
		savedPlaces                        []models.Place
		savedPlanCandidateSet              models.PlanCandidate
		savedPlanCandidateSearchedPlaceIds []string
		expectedPlaces                     []models.Place
	}{
		{
			name:            "find places by plan candidate id",
			planCandidateId: "plan_candidate_id",
			savedPlaces: []models.Place{
				{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
				{Id: "place_id_2", Google: models.GooglePlace{PlaceId: "google_place_id_2"}},
				{Id: "place_id_3", Google: models.GooglePlace{PlaceId: "google_place_id_3"}},
			},
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "plan_candidate_id",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
			savedPlanCandidateSearchedPlaceIds: []string{
				"place_id_1",
				"place_id_2",
			},
			expectedPlaces: []models.Place{
				{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
				{Id: "place_id_2", Google: models.GooglePlace{PlaceId: "google_place_id_2"}},
			},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			defer func(ctx context.Context, db *sql.DB) {
				err := cleanup(ctx, db)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			}(testContext, testDB)

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			// 事前にPlanCandidateSetを保存しておく
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("error while saving plan candidate set: %v", err)
			}

			// 事前にPlanCandidateSearchedPlaceを保存しておく
			for _, searchedPlaceId := range c.savedPlanCandidateSearchedPlaceIds {
				planCandidateSearchedPlaceEntity := generated.PlanCandidateSetSearchedPlace{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: c.savedPlanCandidateSet.Id,
					PlaceID:            searchedPlaceId,
				}
				if err := planCandidateSearchedPlaceEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("error while saving plan candidate searched place: %v", err)
				}
			}

			actualPlaces, err := placeRepository.FindByPlanCandidateId(testContext, c.planCandidateId)
			if err != nil {
				t.Fatalf("error while finding places: %v", err)
			}

			if len(actualPlaces) != len(c.expectedPlaces) {
				t.Fatalf("place expected: %d, actual: %d", len(c.expectedPlaces), len(actualPlaces))
			}
		})
	}
}

func TestPlaceRepository_SaveGooglePlaceDetail(t *testing.T) {
	cases := []struct {
		name              string
		savedPlace        models.Place
		googlePlaceId     string
		googlePlaceDetail models.GooglePlaceDetail
	}{
		{
			name: "save google place detail",
			savedPlace: models.Place{
				Id:     uuid.New().String(),
				Google: models.GooglePlace{PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA"},
			},
			googlePlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
			googlePlaceDetail: models.GooglePlaceDetail{
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
		{
			name: "already saved google place detail",
			savedPlace: models.Place{
				Id: uuid.New().String(),
				Google: models.GooglePlace{
					PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
					PlaceDetail: &models.GooglePlaceDetail{
						OpeningHours: &models.GooglePlaceOpeningHours{
							Periods: []models.GooglePlaceOpeningPeriod{{DayOfWeekOpen: "Monday", DayOfWeekClose: "Monday", OpeningTime: "1030", ClosingTime: "2130"}},
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
						},
					},
				},
			},
			googlePlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
			googlePlaceDetail: models.GooglePlaceDetail{
				OpeningHours: &models.GooglePlaceOpeningHours{
					Periods: []models.GooglePlaceOpeningPeriod{{DayOfWeekOpen: "Monday", DayOfWeekClose: "Monday", OpeningTime: "1030", ClosingTime: "2130"}},
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
				},
			},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			defer func(ctx context.Context, db *sql.DB) {
				err := cleanup(ctx, db)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			}(testContext, testDB)

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, []models.Place{c.savedPlace}); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			if err := placeRepository.SaveGooglePlaceDetail(testContext, c.googlePlaceId, c.googlePlaceDetail); err != nil {
				t.Fatalf("error while saving google place detail: %v", err)
			}

			// GooglePlaceOpeningPeriods が保存されているか確認
			if c.googlePlaceDetail.OpeningHours != nil {
				openingPeriodCount, err := generated.
					GooglePlaceOpeningPeriods(generated.GooglePlaceOpeningPeriodWhere.GooglePlaceID.EQ(c.googlePlaceId)).
					Count(testContext, testDB)
				if err != nil {
					t.Fatalf("error while counting opening periods: %v", err)
				}

				if int(openingPeriodCount) != len(c.googlePlaceDetail.OpeningHours.Periods) {
					t.Fatalf("opening period expected: %d, actual: %d", len(c.googlePlaceDetail.OpeningHours.Periods), openingPeriodCount)
				}
			}

			// GooglePlaceReviews が保存されているか確認
			reviewCount, err := generated.
				GooglePlaceReviews(generated.GooglePlaceReviewWhere.GooglePlaceID.EQ(c.googlePlaceId)).
				Count(testContext, testDB)
			if err != nil {
				t.Fatalf("error while counting reviews: %v", err)
			}

			if int(reviewCount) != len(c.googlePlaceDetail.Reviews) {
				t.Fatalf("review expected: %d, actual: %d", len(c.googlePlaceDetail.Reviews), reviewCount)
			}

			// GooglePhotoReference が保存されているか確認
			for _, photoReference := range c.googlePlaceDetail.PhotoReferences {
				isPhotoReferenceSaved, err := generated.
					GooglePlacePhotoReferences(generated.GooglePlacePhotoReferenceWhere.PhotoReference.EQ(photoReference.PhotoReference)).
					Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("error while checking photo reference existence: %v", err)
				}
				if !isPhotoReferenceSaved {
					t.Fatalf("photo is not saved")
				}
			}
		})
	}
}

func TestPlaceRepository_SaveGooglePlacePhotos(t *testing.T) {
	cases := []struct {
		name              string
		savedPlace        models.Place
		googlePlaceId     string
		googlePlacePhotos []models.GooglePlacePhoto
	}{
		{
			name: "save google place photos with nearby search result",
			savedPlace: models.Place{
				Id: uuid.New().String(),
				Google: models.GooglePlace{
					PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
					PhotoReferences: []models.GooglePlacePhotoReference{
						{
							PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
						},
					},
					Photos: nil,
				},
			},
			googlePlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
			googlePlacePhotos: []models.GooglePlacePhoto{
				{
					PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
					Width:            4032,
					Height:           3024,
					HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					Small: &models.Image{
						Width:  400,
						Height: 400,
						URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
					},
					Large: &models.Image{
						Width:  4032,
						Height: 3024,
						URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
					},
				},
			},
		},
		{
			name: "save google place photos with place detail result",
			savedPlace: models.Place{
				Id: uuid.New().String(),
				Google: models.GooglePlace{
					PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
					PhotoReferences: []models.GooglePlacePhotoReference{
						{
							PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
						},
					},
					PlaceDetail: &models.GooglePlaceDetail{
						PhotoReferences: []models.GooglePlacePhotoReference{
							{
								PhotoReference:   "photo-2-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
								Width:            1920,
								Height:           1080,
								HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
							},
						},
					},
					Photos: nil,
				},
			},
			googlePlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
			googlePlacePhotos: []models.GooglePlacePhoto{
				{
					PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
					Width:            4032,
					Height:           3024,
					HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					Small: &models.Image{
						Width:  400,
						Height: 400,
						URL:    "https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w400-h400",
					},
					Large: &models.Image{
						Width:  4032,
						Height: 3024,
						URL:    "https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w4032-h3024",
					},
				},
				{
					PhotoReference:   "photo-2-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
					Width:            1920,
					Height:           1080,
					HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					Small: &models.Image{
						Width:  400,
						Height: 400,
						URL:    "https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w400-h400",
					},
					Large: &models.Image{
						Width:  4032,
						Height: 3024,
						URL:    "https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w4032-h3024",
					},
				},
			},
		},
		{
			name: "save already saved google place photos",
			savedPlace: models.Place{
				Id: uuid.New().String(),
				Google: models.GooglePlace{
					PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
					PhotoReferences: []models.GooglePlacePhotoReference{
						{
							PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
						},
					},
					Photos: &[]models.GooglePlacePhoto{
						{
							PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
							Small: &models.Image{
								Width:  400,
								Height: 400,
								URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
							},
							Large: &models.Image{
								Width:  4032,
								Height: 3024,
								URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
							},
						},
					},
				},
			},
			googlePlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
			googlePlacePhotos: []models.GooglePlacePhoto{
				{
					PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
					Width:            4032,
					Height:           3024,
					HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					Small: &models.Image{
						Width:  400,
						Height: 400,
						URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
					},
					Large: &models.Image{
						Width:  4032,
						Height: 3024,
						URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
					},
				},
			},
		},
		{
			name: "save duplicated photos",
			savedPlace: models.Place{
				Id: uuid.New().String(),
				Google: models.GooglePlace{
					PlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
					PhotoReferences: []models.GooglePlacePhotoReference{
						{
							PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
						},
					},
					Photos: &[]models.GooglePlacePhoto{
						{
							PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
							Small: &models.Image{
								Width:  400,
								Height: 400,
								URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
							},
							Large: &models.Image{
								Width:  4032,
								Height: 3024,
								URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
							},
						},
						{
							PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
							Width:            4032,
							Height:           3024,
							HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
							Small: &models.Image{
								Width:  400,
								Height: 400,
								URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
							},
							Large: &models.Image{
								Width:  4032,
								Height: 3024,
								URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
							},
						},
					},
				},
			},
			googlePlaceId: "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
			googlePlacePhotos: []models.GooglePlacePhoto{
				{
					PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
					Width:            4032,
					Height:           3024,
					HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
					Small: &models.Image{
						Width:  400,
						Height: 400,
						URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
					},
					Large: &models.Image{
						Width:  4032,
						Height: 3024,
						URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
					},
				},
			},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			defer func(ctx context.Context, db *sql.DB) {
				err := cleanup(ctx, db)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			}(testContext, testDB)

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, []models.Place{c.savedPlace}); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			if err := placeRepository.SaveGooglePlacePhotos(testContext, c.googlePlaceId, c.googlePlacePhotos); err != nil {
				t.Fatalf("error while saving google place photos: %v", err)
			}

			// GooglePlacePhotos が保存されているか確認
			for _, photo := range c.googlePlacePhotos {
				isPhotoSmallSaved, err := generated.
					GooglePlacePhotos(
						generated.GooglePlacePhotoWhere.PhotoReference.EQ(photo.PhotoReference),
						generated.GooglePlacePhotoWhere.URL.EQ(photo.Small.URL),
					).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("error while checking photo existence: %v", err)
				}
				if !isPhotoSmallSaved {
					t.Fatalf("photo small is not saved")
				}

				isPhotoLargeSaved, err := generated.
					GooglePlacePhotos(
						generated.GooglePlacePhotoWhere.PhotoReference.EQ(photo.PhotoReference),
						generated.GooglePlacePhotoWhere.URL.EQ(photo.Large.URL),
					).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("error while checking photo existence: %v", err)
				}
				if !isPhotoLargeSaved {
					t.Fatalf("photo large is not saved")
				}
			}
		})
	}
}

func TestPlaceRepository_SavePlacePhotos(t *testing.T) {
	cases := []struct {
		name               string
		userId             string
		placeId            string
		photoUrl           string
		width              int
		height             int
		preSavedUser       generated.User
		preSavedPlace      generated.Place
		preSavedPlacePhoto generated.PlacePhoto
	}{
		{
			name:     "save place photo",
			userId:   "3b9c288c-3ae6-41be-b375-c5aa6082114d",
			placeId:  "c0bbee6a-acd4-41b6-957e-2aeb83e29d12",
			photoUrl: "https://example.com/photo.jpg",
			width:    1920,
			height:   1080,
			preSavedUser: generated.User{
				ID: "3b9c288c-3ae6-41be-b375-c5aa6082114d",
			},
			preSavedPlace: generated.Place{
				ID: "c0bbee6a-acd4-41b6-957e-2aeb83e29d12",
			},
			preSavedPlacePhoto: generated.PlacePhoto{
				UserID:   "3b9c288c-3ae6-41be-b375-c5aa6082114d",
				PlaceID:  "c0bbee6a-acd4-41b6-957e-2aeb83e29d12",
				PhotoURL: "another-photo.jpg",
			},
		},
		{
			name:     "already saved place photo",
			userId:   "3b9c288c-3ae6-41be-b375-c5aa6082114d",
			placeId:  "c0bbee6a-acd4-41b6-957e-2aeb83e29d12",
			photoUrl: "https://example.com/photo.jpg",
			width:    1920,
			height:   1080,
			preSavedUser: generated.User{
				ID: "3b9c288c-3ae6-41be-b375-c5aa6082114d",
			},
			preSavedPlace: generated.Place{
				ID: "c0bbee6a-acd4-41b6-957e-2aeb83e29d12",
			},
			preSavedPlacePhoto: generated.PlacePhoto{
				UserID:   "3b9c288c-3ae6-41be-b375-c5aa6082114d",
				PlaceID:  "c0bbee6a-acd4-41b6-957e-2aeb83e29d12",
				PhotoURL: "https://example.com/photo.jpg",
			},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}
	testContext := context.Background()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func(ctx context.Context, db *sql.DB) {
				err := cleanup(ctx, db)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			}(testContext, testDB)

			// 事前にPlacePhoto・User・Placeを保存しておく
			if err := c.preSavedUser.Insert(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to insert user: %v", err)
			}
			if err := c.preSavedPlace.Insert(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to insert place: %v", err)
			}
			if err := c.preSavedPlacePhoto.Insert(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to insert place photo: %v", err)
			}

			err := placeRepository.SavePlacePhotos(testContext, c.userId, c.placeId, c.photoUrl, c.width, c.height)
			if err != nil {
				t.Fatalf("error while saving place photo: %v", err)
			}

			_, err = generated.
				PlacePhotos(
					generated.PlacePhotoWhere.UserID.EQ(c.userId),
					generated.PlacePhotoWhere.PlaceID.EQ(c.placeId),
				).Exists(testContext, testDB)
			if err != nil {
				t.Fatalf("error while checking photo existence: %v", err)
			}

		})
	}
}
