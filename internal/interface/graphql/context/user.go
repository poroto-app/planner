package context

import (
	"context"
	"github.com/gin-gonic/gin"
	"poroto.app/poroto/planner/internal/domain/models"
)

const (
	contextAuthUserKey = "auth_user"
)

// SetAuthUser sets the auth user in the context.
func SetAuthUser(c *gin.Context, user *models.User) {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, contextAuthUserKey, user)
	c.Request = c.Request.WithContext(ctx)
}

// GetAuthUser gets the auth user from the context.
func GetAuthUser(ctx context.Context) *models.User {
	if authUser, ok := ctx.Value(contextAuthUserKey).(*models.User); ok {
		return authUser
	}
	return nil
}
