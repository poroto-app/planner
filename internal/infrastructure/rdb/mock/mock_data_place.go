package mock

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

func NewMockGooglePlaceTokyo(placeDetail bool, photos bool) models.GooglePlace {
	googlePlace := models.GooglePlace{
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
	}

	if !placeDetail {
		googlePlace.PlaceDetail = nil
	}

	if !photos {
		googlePlace.Photos = nil
	}

	return googlePlace
}
