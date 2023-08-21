package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionUsers = "users"
)

type UserRepository struct {
	client *firestore.Client
}

func NewUserRepository(ctx context.Context) (*UserRepository, error) {
	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	return &UserRepository{
		client: client,
	}, nil
}

func (u UserRepository) Create(ctx context.Context, user models.User) error {
	if err := u.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		collection := u.collection()

		// 同じFirebaseUIDのユーザーが存在するかどうかを確認する
		alreadyExistsFirebaseUser := false
		docIter := collection.Where("firebase_uid", "==", user.FirebaseUID).Documents(ctx)
		for {
			doc, err := docIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("error while iterating documents: %v", err)
			}

			if doc.Exists() {
				alreadyExistsFirebaseUser = true
				break
			}
		}

		if alreadyExistsFirebaseUser {
			return fmt.Errorf("user already exists")
		}

		// 同じIDのユーザーが存在するかどうかを確認する
		doc := u.doc(user.Id)
		snapshot, err := tx.Get(doc)
		if err != nil {
			return fmt.Errorf("error while getting user: %v", err)
		}

		if snapshot.Exists() {
			return fmt.Errorf("user already exists")
		}

		if _, err := doc.Set(ctx, entity.FromUser(user)); err != nil {
			return fmt.Errorf("error while saving user: %v", err)
		}

		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while creating user: %v", err)
	}

	return nil
}

func (u UserRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*models.User, error) {
	docIter := u.collection().Where("firebase_uid", "==", firebaseUID).Documents(ctx)
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error while iterating documents: %v", err)
		}

		if doc.Exists() {
			var userEntity entity.UserEntity
			if err := doc.DataTo(&userEntity); err != nil {
				return nil, fmt.Errorf("error while converting snapshot to user entity: %v", err)
			}

			user := entity.ToUser(userEntity)
			return &user, nil
		}
	}

	return nil, nil
}

func (u UserRepository) collection() *firestore.CollectionRef {
	return u.client.Collection(collectionUsers)
}

func (u UserRepository) doc(id string) *firestore.DocumentRef {
	return u.collection().Doc(id)
}
