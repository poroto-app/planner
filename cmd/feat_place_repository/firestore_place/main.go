package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/env"
	repo "poroto.app/poroto/planner/internal/infrastructure/firestore"
)

const (
	testGooglePlaceID = "TEST"
)

func init() {
	env.LoadEnv()
}

/**
 * このコードは、仮の Google Place ID を用いて Firestore Place Repository を利用するサンプルコードです。
 *
 * このコードを実行すると、以下の以下の動作が行われます。
 * 1. Google Place ID を `TEST` とした Google Place の仮のデータを Firestore に保存する
 * 2. Place Detail を保存する
 * 3. 画像を保存する
 * 4. 保存したデータを取得する
 * 5. 保存したデータの確認を行う
 * 6. 保存したデータを削除する（CleanUp）
 **/
func main() {
	ctx := context.Background()

	placeRepository, err := repo.NewPlaceRepository(ctx)
	if err != nil {
		log.Fatalf("failed to initialize place repository: %v", err)
	}

	// 検索した場所を保存
	// Google Place ID から他の Place と重複しないように新しいIDをもった Place を作成する
	// もしくは、Google Place ID が既に存在する場合は、既存の Place を返す
	savedPlace, err := placeRepository.SavePlacesFromGooglePlace(ctx, models.GooglePlace{
		PlaceId: testGooglePlaceID,
		Name:    "test",
	})
	if err != nil {
		log.Fatalf("failed to save places from google place: %v", err)
	}
	log.Printf("saved google place[%s] as place[%s]", testGooglePlaceID, savedPlace.Id)

	// Place Detailを保存
	log.Printf("saved google place detail of google place[%s]", testGooglePlaceID)
	if err := placeRepository.SaveGooglePlaceDetail(ctx, testGooglePlaceID, models.GooglePlaceDetail{
		PhotoReferences: []models.GooglePlacePhotoReference{
			{
				PhotoReference: "TEST",
				Width:          100,
				Height:         100,
				HTMLAttributions: []string{
					"TEST",
				},
			},
		},
		Reviews: []models.GooglePlaceReview{
			{
				Rating:                5,
				Text:                  nil,
				Time:                  0,
				AuthorName:            "TEST",
				AuthorUrl:             nil,
				AuthorProfileImageUrl: nil,
				Language:              nil,
				OriginalLanguage:      nil,
			},
		},
		OpeningHours: &models.GooglePlaceOpeningHours{
			Periods: []models.GooglePlaceOpeningPeriod{
				{
					DayOfWeekOpen:  "Monday",
					DayOfWeekClose: "Monday",
					OpeningTime:    "0000",
					ClosingTime:    "0000",
				},
			},
		},
	}); err != nil {
		log.Fatalf("failed to save google place detail: %v", err)
	}

	// 画像を保存
	log.Printf("saved google place photos of google place[%s]", testGooglePlaceID)
	if err := placeRepository.SaveGooglePlacePhotos(ctx, testGooglePlaceID, []models.GooglePlacePhoto{
		{
			PhotoReference: "TEST",
			Width:          100,
			Height:         100,
			HTMLAttributions: []string{
				"TEST",
			},
			Small: utils.StrPointer("TEST"),
			Large: utils.StrPointer("TEST"),
		},
	}); err != nil {
		log.Fatalf("failed to save google place photo: %v", err)
	}

	place, err := placeRepository.FindByGooglePlaceID(ctx, testGooglePlaceID)
	if err != nil {
		log.Fatalf("failed to find place by google place id: %v", err)
	}

	log.Printf("found place[%s] by google place id[%s]", place.Id, testGooglePlaceID)
	log.Printf("place: %+v", place)

	// 値の確認
	if place.Google.PlaceDetail == nil {
		log.Fatalf("place detail is not set")
	}
	log.Printf("place detail: %+v", *place.Google.PlaceDetail)

	if place.Google.PlaceDetail.OpeningHours == nil {
		log.Fatalf("opening hours is not set")
	}
	log.Printf("opening hours: %+v", *place.Google.PlaceDetail.OpeningHours)

	if place.Google.Photos == nil {
		log.Fatalf("photos is not set")
	}
	log.Printf("photos: %+v", *place.Google.Photos)

	if err := CleanUp(ctx); err != nil {
		log.Fatalf("failed to clean up: %v", err)
	}
}

func CleanUp(ctx context.Context) error {
	log.Printf("cleaning up places by google place id[%s]", testGooglePlaceID)

	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return fmt.Errorf("error while initializing firestore client: %v", err)
	}

	log.Printf("fetching places by google place id[%s]", testGooglePlaceID)
	collectionPlaces := client.Collection("places")
	snapshotsPlaces, err := collectionPlaces.Where("google_place_id", "==", testGooglePlaceID).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("error while getting places: %v", err)
	}
	log.Printf("fetched %d places", len(snapshotsPlaces))

	for _, snapshotPlace := range snapshotsPlaces {
		log.Printf("deleting place[%s]", snapshotPlace.Ref.ID)
		if err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			subCollectionGooglePlaces := snapshotPlace.Ref.Collection("google_places")
			subCollectionGoogleReviews := snapshotPlace.Ref.Collection("google_reviews")
			subCollectionGooglePhotos := snapshotPlace.Ref.Collection("google_photos")

			log.Printf("fetching google places of place[%s]", snapshotPlace.Ref.ID)
			snapshotGooglePlaces, err := tx.Documents(subCollectionGooglePlaces).GetAll()
			if err != nil {
				return fmt.Errorf("error while getting google places: %v", err)
			}
			log.Printf("fetched %d google places", len(snapshotGooglePlaces))

			log.Printf("fetching google reviews of place[%s]", snapshotPlace.Ref.ID)
			snapshotGoogleReviews, err := tx.Documents(subCollectionGoogleReviews).GetAll()
			if err != nil {
				return fmt.Errorf("error while getting google reviews: %v", err)
			}
			log.Printf("fetched %d google reviews", len(snapshotGoogleReviews))

			log.Printf("fetching google photos of place[%s]", snapshotPlace.Ref.ID)
			snapshotGooglePhotos, err := tx.Documents(subCollectionGooglePhotos).GetAll()
			if err != nil {
				return fmt.Errorf("error while getting google photos: %v", err)
			}
			log.Printf("fetched %d google photos", len(snapshotGooglePhotos))

			for _, snapshotGooglePlace := range snapshotGooglePlaces {
				log.Printf("deleting google place[%s] of place[%s]", snapshotGooglePlace.Ref.ID, snapshotPlace.Ref.ID)
				if err := tx.Delete(snapshotGooglePlace.Ref); err != nil {
					return fmt.Errorf("error while deleting google place: %v", err)
				}
				log.Printf("deleted google place[%s] of place[%s]", snapshotGooglePlace.Ref.ID, snapshotPlace.Ref.ID)
			}

			for _, snapshotGoogleReview := range snapshotGoogleReviews {
				log.Printf("deleting google review[%s] of place[%s]", snapshotGoogleReview.Ref.ID, snapshotPlace.Ref.ID)
				if err := tx.Delete(snapshotGoogleReview.Ref); err != nil {
					return fmt.Errorf("error while deleting google review: %v", err)
				}
				log.Printf("deleted google review[%s] of place[%s]", snapshotGoogleReview.Ref.ID, snapshotPlace.Ref.ID)
			}

			for _, snapshotGooglePhoto := range snapshotGooglePhotos {
				log.Printf("deleting google photo[%s] of place[%s]", snapshotGooglePhoto.Ref.ID, snapshotPlace.Ref.ID)
				if err := tx.Delete(snapshotGooglePhoto.Ref); err != nil {
					return fmt.Errorf("error while deleting google photo: %v", err)
				}
				log.Printf("deleted google photo[%s] of place[%s]", snapshotGooglePhoto.Ref.ID, snapshotPlace.Ref.ID)
			}

			log.Printf("deleting place[%s]", snapshotPlace.Ref.ID)
			if err := tx.Delete(snapshotPlace.Ref); err != nil {
				return fmt.Errorf("error while deleting place: %v", err)
			}
			log.Printf("deleted place[%s]", snapshotPlace.Ref.ID)

			return nil
		}); err != nil {
			return fmt.Errorf("error while running transaction: %v", err)
		}
	}

	return nil
}
