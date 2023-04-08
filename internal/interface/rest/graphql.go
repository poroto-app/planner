package rest

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"poroto.app/poroto/planner/graphql"
	"poroto.app/poroto/planner/graphql/generated"
)

func GraphQlPlayGround(c *gin.Context) {
	h := playground.Handler("GraphQL", "/graphql")
	h.ServeHTTP(c.Writer, c.Request)
}

func GraphQlQuery(c *gin.Context) {
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: &graphql.Resolver{}})
	h := handler.NewDefaultServer(schema)
	h.ServeHTTP(c.Writer, c.Request)
}
