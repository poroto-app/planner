package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

func init() {
	env := os.Getenv("ENV")
	if "" == env {
		env = "development"
	}

	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalf("error while loading .env.local: %v", err)
	}

	if err := godotenv.Load(".env." + env); err != nil {
		log.Fatalf("error while loading .env.%s: %v", env, err)
	}
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

	if err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		log.Println("Start transaction")

		collection := client.Collection("plans")

		var updatedPlans []entity.PlanEntity
		documentIter := tx.DocumentRefs(collection)
		for {
			doc, err := documentIter.Next()
			if err != nil {
				break
			}

			planDoc, err := tx.Get(doc)
			if err != nil {
				return err
			}

			log.Printf("Updating plan: %s\n", doc.ID)

			var plan entity.PlanEntity
			if err := planDoc.DataTo(&plan); err != nil {
				return err
			}

			for i, place := range plan.Places {
				if place.Images != nil && len(place.Images) > 0 {
					log.Printf("Skip place: %s\n", place.Id)
					continue
				}

				// images が保存されていない場合は、上書きする
				place.Images = make([]entity.ImageEntity, len(place.Photos))
				for i, photo := range place.Photos {
					place.Images[i] = entity.ImageEntity{
						Small: nil,
						Large: &photo,
					}
				}

				plan.Places[i] = place
			}

			updatedPlans = append(updatedPlans, plan)
		}

		for _, plan := range updatedPlans {
			doc := collection.Doc(plan.Id)
			if err := tx.Update(doc, []firestore.Update{
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
		}

		return nil
	}); err != nil {
		log.Fatalf("error while transaction: %v", err)
	}
}
