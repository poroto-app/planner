package rdb

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"testing"
)

func TestUserRepository_UpdateProfile(t *testing.T) {
	cases := []struct {
		name       string
		savedUsers generated.UserSlice
		userId     string
		userName   *string
		photoUrl   *string
		expected   generated.User
	}{
		{
			name: "update user profile with empty values",
			savedUsers: generated.UserSlice{
				{
					ID:       "user-1",
					Name:     null.StringFrom("test-user"),
					PhotoURL: null.StringFrom("https://example.com/test-user.jpg"),
				},
			},
			userId:   "user-1",
			userName: nil,
			photoUrl: nil,
			expected: generated.User{
				ID:       "user-1",
				Name:     null.StringFrom("test-user"),
				PhotoURL: null.StringFrom("https://example.com/test-user.jpg"),
			},
		},
		{
			name: "update user profile with new values",
			savedUsers: generated.UserSlice{
				{
					ID:       "user-1",
					Name:     null.StringFrom("test-user"),
					PhotoURL: null.StringFrom("https://example.com/test-user.jpg"),
				},
			},
			userId:   "user-1",
			userName: utils.ToPointer("new-user"),
			photoUrl: utils.ToPointer("https://example.com/new-user.jpg"),
			expected: generated.User{
				ID:       "user-1",
				Name:     null.StringFrom("new-user"),
				PhotoURL: null.StringFrom("https://example.com/new-user.jpg"),
			},
		},
	}

	userRepository, err := NewUserRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup(testContext, testDB); err != nil {
					t.Errorf("error cleaning up: %v", err)
				}
			})

			// データを準備
			if _, err := c.savedUsers.InsertAll(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to insert user: %v", err)
			}

			// テスト対象の関数を実行
			if err := userRepository.UpdateProfile(testContext, c.userId, c.userName, c.photoUrl); err != nil {
				t.Fatalf("failed to update user profile: %v", err)
			}

			updatedUser, err := generated.Users(generated.UserWhere.ID.EQ(c.userId)).One(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to find user: %v", err)
			}

			if diff := cmp.Diff(
				c.expected,
				*updatedUser,
				cmpopts.IgnoreFields(generated.User{}, "CreatedAt", "UpdatedAt"),
			); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}
