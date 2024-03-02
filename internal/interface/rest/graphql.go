package rest

import (
	"database/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"poroto.app/poroto/planner/internal/domain/services/user"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/interface/graphql/generated"
	"poroto.app/poroto/planner/internal/interface/graphql/resolver"
)

func GraphQlPlayGround(c *gin.Context) {
	h := playground.Handler("GraphQL", "/graphql")
	h.ServeHTTP(c.Writer, c.Request)
}

func GraphQlQueryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger, err := utils.NewLogger(utils.LoggerOption{
			Tag: "GraphQL",
		})
		if err != nil {
			log.Println("error while initializing Logger: ", err)
			c.JSON(500, gin.H{
				"error": "internal server error",
			})
		}

		userService, err := user.NewService(c.Request.Context(), db)
		if err != nil {
			logger.Error("error while initializing user service", zap.Error(err))
			c.JSON(500, gin.H{
				"error": "internal server error",
			})
		}

		schema := generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
			Logger:      logger,
			DB:          db,
			UserService: userService,
		}})
		h := handler.NewDefaultServer(schema)
		h.ServeHTTP(c.Writer, c.Request)
	}
}
