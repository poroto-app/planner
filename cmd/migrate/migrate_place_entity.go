package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/option"
	"log"
	"os"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/env"
	"time"
)

func init() {
	env.LoadEnv()
}

// PlanEntity 2023/10/4
type PlanEntity struct {
	Id            string               `firestore:"id"`
	Name          string               `firestore:"name"`
	Places        []PlaceEntity        `firestore:"places"`
	GeoHash       *string              `firestore:"geohash,omitempty"`
	TimeInMinutes int                  `firestore:"time_in_minutes"`
	Transitions   *[]TransitionsEntity `firestore:"transitions,omitempty"`
	CreatedAt     time.Time            `firestore:"created_at,omitempty,serverTimestamp"`
	UpdatedAt     time.Time            `firestore:"updated_at,omitempty"`
	AuthorId      *string              `firestore:"author_id,omitempty"`
}

// PlaceEntity 2023/10/4
type PlaceEntity struct {
	Id                    string                     `firestore:"id"`
	GooglePlaceId         *string                    `firestore:"google_place_id"`
	Name                  string                     `firestore:"name"`
	Location              GeoLocationEntity          `firestore:"location"`
	Thumbnail             *string                    `firestore:"thumbnail"`
	Photos                []string                   `firestore:"photos"`
	Images                []ImageEntity              `firestore:"images"`
	EstimatedStayDuration int                        `firestore:"estimated_stay_duration"`
	GooglePlaceReviews    *[]GooglePlaceReviewEntity `firestore:"google_place_reviews,omitempty"`
}

// TransitionsEntity 2023/10/4
type TransitionsEntity struct {
	FromPlaceId *string `firestore:"from,omitempty"`
	ToPlaceId   string  `firestore:"to"`
	Duration    int     `firestore:"duration"`
}

// GeoLocationEntity 2023/10/4
type GeoLocationEntity struct {
	Latitude  float64 `firestore:"latitude"`
	Longitude float64 `firestore:"longitude"`
}

// ImageEntity 2023/10/4
type ImageEntity struct {
	Small *string `firestore:"small,omitempty"`
	Large *string `firestore:"large,omitempty"`
}

// GooglePlaceReviewEntity 2023/10/4
type GooglePlaceReviewEntity struct {
	Rating           int     `firestore:"rating"`
	Text             *string `firestore:"text,omitempty"`
	Time             int     `firestore:"time"`
	AuthorName       string  `firestore:"author_name"`
	AuthorUrl        *string `firestore:"author_url,omitempty"`
	AuthorProfileUrl *string `firestore:"author_profile_url,omitempty"`
	Language         *string `firestore:"language,omitempty"`
	OriginalLanguage *string `firestore:"original_language,omitempty"`
}

// 保存されたプランに含まれる場所の Photos, Thumbnail
// プロパティを削除し、Imagesに書き換える
func main() {
	log.Println("Start migration")
	log.Printf("Environment: %s\n", os.Getenv("ENV"))

	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		log.Fatalf("error while initializing firestore client: %v", err)
	}

	collection := client.Collection("plans")
	documentIter := collection.Documents(ctx)
	for {
		snapshot, err := documentIter.Next()
		if err != nil {
			break
		}

		if snapshot.Ref == nil {
			continue
		}

		if err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			log.Println("Start transaction====================")
			log.Printf("Updating plan: %s\n", snapshot.Ref.ID)

			var plan PlanEntity
			if err := snapshot.DataTo(&plan); err != nil {
				return err
			}

			for i, place := range plan.Places {
				if place.Images != nil && len(place.Images) > 0 {
					log.Printf("Skip place: %s\n", place.Id)
					continue
				}

				// images が保存されていない場合は、上書きする
				place.Images = make([]ImageEntity, len(place.Photos))
				for i, photo := range place.Photos {
					place.Images[i] = ImageEntity{
						Small: nil,
						Large: utils.StrPointer(photo),
					}
				}

				plan.Places[i] = place
			}

			if err := tx.Update(collection.Doc(snapshot.Ref.ID), []firestore.Update{
				{
					Path:  "places",
					Value: plan.Places,
				},
				{
					Path:  "updated_at",
					Value: firestore.ServerTimestamp,
				},
			}); err != nil {
				return err
			}

			log.Printf("Updated plan: %s", plan.Id)
			log.Printf("End transaction====================")

			return nil
		}); err != nil {
			log.Fatalf("error while transaction: %v", err)
		}
	}

}
