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

// Verify firebaseUid と tokenId　から取得されるユーザーが同一であるかを確認する。
func (f *FirebaseAuth) Verify(
	ctx context.Context,
	firebaseUid string,
	tokenId string,
) (bool, error) {
	token, err := f.client.VerifyIDToken(ctx, tokenId)
	if err != nil {
		return false, fmt.Errorf("error while verifying firebase token: %v", err)
	}

	if token.UID != firebaseUid {
		return false, nil
	}

	return true, nil
}

func (f *FirebaseAuth) GetUser(ctx context.Context, firebaseUid string) (*auth.UserRecord, error) {
	return f.client.GetUser(ctx, firebaseUid)
}
