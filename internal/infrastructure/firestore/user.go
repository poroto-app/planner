package firestore

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Printf("[%s] Start checking if the user with same firebase id already exists\n", user.Id)
		alreadyExistsFirebaseUser := false
		firebaseUserQuery := collection.Where("firebase_uid", "==", user.FirebaseUID)
		docIter := tx.Documents(firebaseUserQuery)
		for {
			doc, err := docIter.Next()
			if err == iterator.Done {
				// 検索結果なし
				break
			} else if err != nil {
				return fmt.Errorf("error while iterating documents: %v", err)
			} else if doc.Exists() {
				// すでに同じFirebaseUIDのデータが存在する
				alreadyExistsFirebaseUser = true
				break
			}
		}

		if alreadyExistsFirebaseUser {
			return fmt.Errorf("user already exists")
		}

		// 同じIDのユーザーが存在するかどうかを確認する
		log.Printf("[%s] Start checking if the user with same id already exists\n", user.Id)
		doc := u.doc(user.Id)
		snapshot, err := tx.Get(doc)
		if err != nil && status.Code(err) != codes.NotFound {
			return fmt.Errorf("error while getting user: %v", err)
		} else if snapshot.Exists() {
			return fmt.Errorf("user already exists")
		}

		// ユーザーデータ書き込み
		log.Printf("[%s] Start saving user\n", user.Id)
		if err := tx.Set(doc, entity.FromUser(user)); err != nil {
			return fmt.Errorf("error while saving user: %v", err)
		}

		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while creating user: %v", err)
	}

	return nil
}

func (u UserRepository) Find(ctx context.Context, id string) (*models.User, error) {
	doc := u.doc(id)
	snapshot, err := doc.Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, fmt.Errorf("error while getting user: %v", err)
	}

	if !snapshot.Exists() {
		return nil, nil
	}

	var userEntity entity.UserEntity
	if err := snapshot.DataTo(&userEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to user entity: %v", err)
	}

	user := entity.ToUser(userEntity)
	return &user, nil
}

func (u UserRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*models.User, error) {
	log.Printf("Start finding user by firebase uid: %s\n", firebaseUID)
	docIter := u.collection().Where("firebase_uid", "==", firebaseUID).Documents(ctx)
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			log.Printf("User not found by firebase uid: %s\n", firebaseUID)
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error while iterating documents: %v", err)
		}

		log.Printf("User found by firebase uid: %s\n", firebaseUID)
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
