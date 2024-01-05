package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func init() {
	os.Setenv("ENV", "development")
	env.LoadEnv()
}

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&loc=%s&tls=%v&interpolateParams=%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		"Asia%2FTokyo",
		os.Getenv("ENV") != "development",
		true,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	boil.SetDB(db)
	boil.DebugMode = true

	log.Println(dsn)

	cleanup(context.Background(), db)
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}
	log.Println("ping ok")

	repository, err := rdb.NewPlaceRepository(db)
	if err != nil {
		panic(err)
	}

	testGooglePlace := models.GooglePlace{
		PlaceId:  "ChIJ7WoyEQr9GGAREzlMT6J-JhA",
		Name:     "ブルーム",
		Types:    []string{"cafe", "food", "point_of_interest", "store", "establishment"},
		Location: models.GeoLocation{Latitude: 35.5718006, Longitude: 139.3898712},
		PhotoReferences: []models.GooglePlacePhotoReference{
			{
				PhotoReference: "AWU5eFjiROQJEeMpt7Hh2Pv_YhnQmQ3g-bGIEYxQJqVuj8R61cW-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
				Width:          4032,
				Height:         3024,
				HTMLAttributions: []string{
					"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>",
				},
			},
		},
		PriceLevel:       0,
		Rating:           3.9000000953674316,
		UserRatingsTotal: 38,
		Vicinity:         utils.StrPointer("日本、〒150-0001 東京都渋谷区神宮前５丁目２８−１"),
		Photos:           nil,
		PlaceDetail:      nil,
	}

	placeDetail := models.GooglePlaceDetail{
		OpeningHours: &models.GooglePlaceOpeningHours{
			Periods: []models.GooglePlaceOpeningPeriod{
				{
					DayOfWeekOpen:  "Monday",
					DayOfWeekClose: "Monday",
					OpeningTime:    "1030",
					ClosingTime:    "2130",
				},
				{
					DayOfWeekOpen:  "Tuesday",
					DayOfWeekClose: "Tuesday",
					OpeningTime:    "1030",
					ClosingTime:    "2130",
				},
				{
					DayOfWeekOpen:  "Wednesday",
					DayOfWeekClose: "Wednesday",
					OpeningTime:    "1030",
					ClosingTime:    "2130",
				},
				{
					DayOfWeekOpen:  "Thursday",
					DayOfWeekClose: "Thursday",
					OpeningTime:    "1030",
					ClosingTime:    "2130",
				},
				{
					DayOfWeekOpen:  "Friday",
					DayOfWeekClose: "Friday",
					OpeningTime:    "1030",
					ClosingTime:    "2130",
				},
				{
					DayOfWeekOpen:  "Saturday",
					DayOfWeekClose: "Saturday",
					OpeningTime:    "1030",
					ClosingTime:    "2130",
				},
			},
		},
		PhotoReferences: []models.GooglePlacePhotoReference{
			{
				PhotoReference: "AWU5eFgYAi-hoRwbOOkzvamLTtNi-lHUN-8Cbcl2xGP49EwZ5xzfo10jvcvuegwztrqV1iJmAjtG0XVs8Ph52lfav7mROP2Srh7h74CMNtXsQBKhIdFsjLp03zOcpfAWNkHqi4H54hyJ3VekpHvbiWOrayPbhnmWchlB5sLwcn17snJQ2uWA",
				Width:          4032,
				Height:         3024,
				HTMLAttributions: []string{
					"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>",
				},
			},
			{
				PhotoReference: "AWU5eFgjBov_JHzMF3w_3rqW0KQK7bmnDzJx7ol0hHW3JPItzw0Ig15zSW8BTABC6mbLJKbwtPgi_56kYDU2qVlU39X92FzSG1bX6pXO9e-tio_BLY3I5p1bcsNkwS7DYluewdPD9_fnip8AMM8zBoF5s2ZY_0GnjyUjUT7TpwmEHdYuynUA",
				Width:          4032,
				Height:         3024,
				HTMLAttributions: []string{
					"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>",
				},
			},
			{
				PhotoReference: "AWU5eFh7d3vAnP6K1mA9qPNwkhrkvLyBoFPfPJ0aNAeIvd8w9Kwvbb5TSWY51ZYkMLndj5OWfUgGQ4SsHMiLZiKOdJ5QOppGFTUJ0ZlI0CIvvQ1j2aVry4JxGagskidptsJuSD2cnfBuLBGFka_2CpNRHo_bIrYi8UGCke16vIszKpn1D3_j",
				Height:         5333,
				Width:          3000,
			},
		},
		Reviews: []models.GooglePlaceReview{
			{
				Rating:                4,
				Text:                  utils.StrPointer("トーストとコーヒーのセットを頼んだら、フルーツも着いてて嬉しかったです。"),
				Time:                  1648126226,
				AuthorName:            "Miho Toyomura",
				AuthorProfileImageUrl: utils.StrPointer("https://lh3.googleusercontent.com/a/ACg8ocKaPr9FWIiqs88c_qhl11GF5KF3F3zQn4XDeZCJvkrc=s128-c0x00000000-cc-rp-mo-ba5"),
				AuthorUrl:             utils.StrPointer("https://www.google.com/maps/contrib/117007323051675636986/reviews"),
			},
			{
				Rating:                5,
				Text:                  utils.StrPointer("ハンバーグは繋ぎが少なく(使ってない？)しっかりと肉を感じるけれどパサついておらずふっくらとジューシーで食べ応えがありました。ライスを食べる時はまさしくこういうハンバーグが理想的だと思います。とんかつも、衣に工夫がされており肉も衣も厚くジューシー。火の入り具合もちょうもよく、ムチっと柔らかな食感でこちらも絶品でした。これでコーヒーついて800円は安すぎる・・・！私は喫煙しませんが特にタバコは気にしないタイプなので、吸う人だったらさらに天国だろうなと。エビチリパスタなどシェフの創意工夫が凝らされたメニューが他にもあったので必ずまた行きます。"),
				Time:                  1618085426,
				AuthorName:            "official tommy",
				AuthorProfileImageUrl: utils.StrPointer("https://lh3.googleusercontent.com/a-/ALV-UjUg1pWpb0FEWmd_wD8wQ5y5NPqCU7qZM9rnp00GHZYagec=s128-c0x00000000-cc-rp-mo-ba4"),
				AuthorUrl:             utils.StrPointer("https://www.google.com/maps/contrib/100755868001879781001/reviews"),
			},
		},
	}

	photos := []models.GooglePlacePhoto{
		{
			PhotoReference: "AWU5eFgYAi-hoRwbOOkzvamLTtNi-lHUN-8Cbcl2xGP49EwZ5xzfo10jvcvuegwztrqV1iJmAjtG0XVs8Ph52lfav7mROP2Srh7h74CMNtXsQBKhIdFsjLp03zOcpfAWNkHqi4H54hyJ3VekpHvbiWOrayPbhnmWchlB5sLwcn17snJQ2uWA",
			Width:          4032,
			Height:         3024,
			Small:          utils.StrPointer("https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w400-h400"),
			Large:          utils.StrPointer("https://lh3.googleusercontent.com/places/ANXAkqEs-dl0rT1eITFJ3j4kMuMKgoRtb-Ws8lhKidWPL7LU4e-57yzhuN5UisB2S-fn4yj23gDQIrlQReGkuMI1Y8QU3ZsxQk2wwgw=s1600-w400-h400"),
			HTMLAttributions: []string{
				"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>",
			},
		},
		{
			PhotoReference: "AWU5eFgjBov_JHzMF3w_3rqW0KQK7bmnDzJx7ol0hHW3JPItzw0Ig15zSW8BTABC6mbLJKbwtPgi_56kYDU2qVlU39X92FzSG1bX6pXO9e-tio_BLY3I5p1bcsNkwS7DYluewdPD9_fnip8AMM8zBoF5s2ZY_0GnjyUjUT7TpwmEHdYuynUA",
			Width:          4032,
			Height:         3024,
			Small:          utils.StrPointer("https://lh3.googleusercontent.com/places/ANXAkqFeWqaRRrXTpjDZrwEI4g7_ui5Xd_3be9W7IiTN4ATrRTTS5Ij83QYtcKseL8v8T9irGwdILcA2MciaTXQN_rSsDe5X0TV_ttA=s1600-w1000-h1000"),
			Large:          utils.StrPointer("https://lh3.googleusercontent.com/places/ANXAkqFeWqaRRrXTpjDZrwEI4g7_ui5Xd_3be9W7IiTN4ATrRTTS5Ij83QYtcKseL8v8T9irGwdILcA2MciaTXQN_rSsDe5X0TV_ttA=s1600-w1000-h1000"),
			HTMLAttributions: []string{
				"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>",
			},
		},
		{
			PhotoReference: "AWU5eFjiROQJEeMpt7Hh2Pv_YhnQmQ3g-bGIEYxQJqVuj8R61cW-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
			Width:          4032,
			Height:         3024,
			HTMLAttributions: []string{
				"<a href=\"https://maps.google.com/maps/contrib/100755868001879781001\">A Google User</a>",
			},
			Small: utils.StrPointer("https://lh3.googleusercontent.com/places/ANXAkqFeWqaRRrXTpjDZrwEI4g7_ui5Xd_3be9W7IiTN4ATrRTTS5Ij83QYtcKseL8v8T9irGwdILcA2MciaTXQN_rSsDe5X0TV_ttA=s1600-w1000-h1000"),
			Large: utils.StrPointer("https://lh3.googleusercontent.com/places/ANXAkqFeWqaRRrXTpjDZrwEI4g7_ui5Xd_3be9W7IiTN4ATrRTTS5Ij83QYtcKseL8v8T9irGwdILcA2MciaTXQN_rSsDe5X0TV_ttA=s1600-w1000-h1000"),
		},
	}

	place, err := repository.SavePlacesFromGooglePlace(context.Background(), testGooglePlace)
	if err != nil {
		panic(err)
	}

	if diff := cmp.Diff(testGooglePlace, place.Google); diff != "" {
		log.Printf("+want, -got:\n%s", diff)
	}

	// PlaceDetailを保存
	if err := repository.SaveGooglePlaceDetail(context.Background(), testGooglePlace.PlaceId, placeDetail); err != nil {
		panic(err)
	}

	testGooglePlace.PlaceDetail = &placeDetail

	if err := repository.SaveGooglePlacePhotos(context.Background(), testGooglePlace.PlaceId, photos); err != nil {
		panic(err)
	}

	testGooglePlace.Photos = &photos

	placeByGooglePlaceId, err := repository.FindByGooglePlaceID(context.Background(), testGooglePlace.PlaceId)
	if err != nil {
		panic(err)
	}

	if diff := cmp.Diff(testGooglePlace, placeByGooglePlaceId.Google); diff != "" {
		log.Printf("+want, -got:\n%s", diff)
	}

	placesByLocation, err := repository.FindByLocation(context.Background(), testGooglePlace.Location)
	if err != nil {
		panic(err)
	}

	log.Println(placesByLocation)
}

func cleanup(ctx context.Context, db *sql.DB) {
	// Truncate all tables
	tx, _ := db.Begin()
	_, _ = generated.GooglePlaceTypes().DeleteAll(ctx, tx)
	_, _ = generated.GooglePlacePhotoAttributions().DeleteAll(ctx, tx)
	_, _ = generated.GooglePlacePhotos().DeleteAll(ctx, tx)
	_, _ = generated.GooglePlacePhotoReferences().DeleteAll(ctx, tx)
	_, _ = generated.GooglePlaceReviews().DeleteAll(ctx, tx)
	_, _ = generated.GooglePlaceOpeningPeriods().DeleteAll(ctx, tx)
	_, _ = generated.GooglePlaces().DeleteAll(ctx, tx)
	_, _ = generated.Places().DeleteAll(ctx, tx)
	_ = tx.Commit()
}
