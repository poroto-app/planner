package rest

import (
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/domain/services/plan"
	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
	"poroto.app/poroto/planner/internal/domain/services/plangen"
	"poroto.app/poroto/planner/internal/domain/services/user"
	"poroto.app/poroto/planner/internal/domain/utils"
	gcontext "poroto.app/poroto/planner/internal/interface/graphql/context"
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

		planService, err := plan.NewService(c.Request.Context(), db)
		if err != nil {
			logger.Error("error while initializing plan service", zap.Error(err))
			c.JSON(500, gin.H{
				"error": "internal server error",
			})
		}

		planGenService, err := plangen.NewService(db)
		if err != nil {
			logger.Error("error while initializing plan gen service", zap.Error(err))
			c.JSON(500, gin.H{
				"error": "internal server error",
			})
		}

		planCandidateService, err := plancandidate.NewService(c.Request.Context(), db)
		if err != nil {
			logger.Error("error while initializing plan candidate service", zap.Error(err))
			c.JSON(500, gin.H{
				"error": "internal server error",
			})
		}

		placeService, err := place.NewService(c.Request.Context(), db)
		if err != nil {
			logger.Error("error while initializing place service", zap.Error(err))
			c.JSON(500, gin.H{
				"error": "internal server error",
			})
		}

		schema := generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
			Logger:               logger,
			DB:                   db,
			UserService:          userService,
			PlanService:          planService,
			PlanCandidateService: planCandidateService,
			PlanGenService:       planGenService,
			PlaceService:         placeService,
		}})
		h := handler.NewDefaultServer(schema)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// GraphqlAuthMiddleware Authorization Header が設定されている場合のみ
// 対応するユーザーを取得し、contextにセットする
func (s Server) GraphqlAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.logger.Info("GraphqlAuthMiddleware")

		var idToken string
		_, err := fmt.Sscanf(c.GetHeader("Authorization"), "Bearer %s", &idToken)
		if err != nil || idToken == "" {
			s.logger.Debug("no authorization header")
			c.Next()
			return
		}

		firebaseUid, err := s.firebaseAuth.GetFirebaseUIDFromTokenId(c.Request.Context(), idToken)
		if err != nil {
			s.logger.Warn(
				"error while getting firebase uid from token id",
				zap.Error(err),
			)
			c.Next()
			return
		}

		user, err := s.userRepository.FindByFirebaseUID(c.Request.Context(), *firebaseUid)
		if err != nil {
			s.logger.Warn(
				"error while getting user by firebase uid",
				zap.Error(err),
			)
			c.Next()
			return
		}

		s.logger.Debug(
			"set auth user to context",
			zap.String("userId", user.Id),
		)
		gcontext.SetAuthUser(c, user)
		c.Next()
	}
}
