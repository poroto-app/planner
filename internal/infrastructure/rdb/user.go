package rdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	return &UserRepository{
		db: db,
	}, nil
}

func (u UserRepository) GetDB() *sql.DB {
	return u.db
}

func (u UserRepository) Create(ctx context.Context, user models.User) error {
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error while starting transaction: %v", err)
	}

	// 同じFirebaseUIDのユーザーが存在するかどうかを確認する
	exists, err := generated.Users(generated.UserWhere.FirebaseUID.EQ(user.FirebaseUID)).Exists(ctx, tx)
	if err != nil {
		return fmt.Errorf("error while checking if the user with same firebase id already exists: %v", err)
	}
	if exists {
		return fmt.Errorf("user with same firebase id already exists")
	}

	// ユーザーを作成する
	userEntity := factory.NewUserEntityFromUser(user)
	if err := userEntity.Insert(ctx, tx, boil.Infer()); err != nil {
		return fmt.Errorf("error while inserting user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error while committing transaction: %v", err)
	}

	return nil
}

func (u UserRepository) Find(ctx context.Context, id string) (*models.User, error) {
	userEntity, err := generated.FindUser(ctx, u.db, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error while finding user: %v", err)
	}

	if userEntity == nil {
		return nil, nil
	}

	return factory.NewUserFromUserEntity(*userEntity), nil
}

func (u UserRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*models.User, error) {
	userEntity, err := generated.Users(generated.UserWhere.FirebaseUID.EQ(firebaseUID)).One(ctx, u.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error while finding user: %v", err)
	}

	if userEntity == nil {
		return nil, nil
	}

	return factory.NewUserFromUserEntity(*userEntity), nil
}

func (u UserRepository) UpdateProfile(ctx context.Context, userId string, name *string, photoUrl *string) error {
	if err := runTransaction(ctx, u, func(ctx context.Context, tx *sql.Tx) error {
		userEntity, err := generated.FindUser(ctx, tx, userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("user not found")
			}
			return fmt.Errorf("error while finding user: %v", err)
		}

		if name != nil {
			userEntity.Name = null.StringFromPtr(name)
		}
		if photoUrl != nil {
			userEntity.PhotoURL = null.StringFromPtr(photoUrl)
		}

		if _, err := userEntity.Update(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("error while updating user: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while updating user profile: %v", err)
	}

	return nil
}
