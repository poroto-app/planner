//go:generate go run github.com/99designs/gqlgen

package resolver

import (
	"database/sql"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/services/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Logger      *zap.Logger
	DB          *sql.DB
	UserService *user.Service
}
