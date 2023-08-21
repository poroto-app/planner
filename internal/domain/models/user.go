package models

type User struct {
	Id          string
	FirebaseUID string
	Name        string
	Email       *string
	PhotoUrl    *string
}
