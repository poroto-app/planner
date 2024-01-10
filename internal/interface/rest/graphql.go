package rest

import (
	"database/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"poroto.app/poroto/planner/internal/interface/graphql/generated"
	"poroto.app/poroto/planner/internal/interface/graphql/resolver"
)

func GraphQlPlayGround(c *gin.Context) {
	h := playground.Handler("GraphQL", "/graphql")
	h.ServeHTTP(c.Writer, c.Request)
}

func GraphQlQueryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		schema := generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
			DB: db,
		}})
		h := handler.NewDefaultServer(schema)
		h.ServeHTTP(c.Writer, c.Request)
	}
}
