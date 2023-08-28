package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"time"
)

type UserEntity struct {
	Id          string    `firestore:"id"`
	FirebaseUID string    `firestore:"firebase_uid"`
	Name        string    `firestore:"name"`
	Email       *string   `firestore:"email,omitempty"`
	PhotoUrl    *string   `firestore:"photo_url,omitempty"`
	CreatedAt   time.Time `firestore:"created_at,omitempty,serverTimestamp"`
	UpdatedAt   time.Time `firestore:"updated_at,omitempty"`
}

func FromUser(user models.User) UserEntity {
	return UserEntity{
		Id:          user.Id,
		FirebaseUID: user.FirebaseUID,
		Name:        user.Name,
		Email:       user.Email,
		PhotoUrl:    user.PhotoUrl,
		UpdatedAt:   time.Now(),
	}
}

func ToUser(entity UserEntity) models.User {
	return models.User{
		Id:          entity.Id,
		FirebaseUID: entity.FirebaseUID,
		Name:        entity.Name,
		Email:       entity.Email,
		PhotoUrl:    entity.PhotoUrl,
	}
}
