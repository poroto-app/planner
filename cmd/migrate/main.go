package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
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

// TODO: DELETE ME
// Firestoreに保存されるPlaceオブジェクトにIDを追加するためのスクリプト
func main() {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), option.WithCredentialsFile("secrets/google-credential.json"))
	if err != nil {
		log.Fatalf("error while initializing firestore client: %v", err)
	}

	//　既にIDが付与されていたとしても、新しいIDを生成し、上書きする
	//　（既に設定されているIDはPlaces APIによって取得されるIDのため）
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		query := tx.Documents(client.Collection("plans"))
		snapshots, err := query.GetAll()
		if err != nil {
			return fmt.Errorf("error while getting snapshots: %v", err)
		}

		for _, snapshot := range snapshots {
			log.Printf("Updating plan: %s\n", snapshot.Ref.ID)
			if err != nil {
				return fmt.Errorf("error while getting snapshot: %v", err)
			}

			var planEntity entity.PlanEntity
			if err = snapshot.DataTo(&planEntity); err != nil {
				return fmt.Errorf("error while converting snapshot to plan entity: %v", err)
			}

			for i := 0; i < len(planEntity.Places); i++ {
				planEntity.Places[i].Id = uuid.New().String()
			}

			if err := tx.Update(snapshot.Ref, []firestore.Update{
				{
					Path:  "places",
					Value: planEntity.Places,
				},
			}); err != nil {
				return fmt.Errorf("error while saving plan: %v", err)
			}

			log.Printf("Updated plan: %s\n", snapshot.Ref.ID)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("error while updating plans: %v", err)
	}
}
