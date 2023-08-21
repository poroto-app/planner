package auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"

	firebase "firebase.google.com/go/v4"
)

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(ctx context.Context) (*FirebaseAuth, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firebase app: %v", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firebase auth client: %v", err)
	}

	return &FirebaseAuth{
		client: client,
	}, nil
}
