package rest

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"os"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/auth"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	port           string
	mode           string
	firebaseAuth   auth.FirebaseAuth
	userRepository repository.UserRepository
	logger         zap.Logger
}

const (
	ServerModeDevelopment = "development"
	ServerModeStaging     = "staging"
	ServerModeProduction  = "production"
)

func NewRestServer(ctx context.Context, db *sql.DB, env string) (*Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "RestServer",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing Logger: %w", err)
	}

	firebaseAuth, err := auth.NewFirebaseAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firebase auth: %w", err)
	}

	userRepository, err := rdb.NewUserRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing user repository: %w", err)
	}

	return &Server{
		port:           port,
		mode:           serverModeFromEnv(env),
		firebaseAuth:   *firebaseAuth,
		userRepository: userRepository,
		logger:         *logger,
	}, nil
}

func (s Server) ServeHTTP(db *sql.DB) error {
	if s.isStaging() || s.isProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST"},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		AllowOriginFunc: func(origin string) bool {
			if s.isDevelopment() {
				return true
			}

			u, err := url.Parse(origin)
			if err != nil {
				return false
			}

			protocol := os.Getenv("WEB_PROTOCOL")
			host := os.Getenv("WEB_HOST")
			return u.Scheme == protocol && u.Host == host
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from planner API",
		})
	})

	groupGraphql := r.Group("/graphql")
	{
		groupGraphql.Use(s.GraphqlAuthMiddleware())
		groupGraphql.POST("", GraphQlQueryHandler(db))
		if s.isDevelopment() || s.isStaging() {
			groupGraphql.GET("/playground", GraphQlPlayGround)
		}
	}

	if err := r.Run(":" + s.port); err != nil {
		return nil
	}

	return nil
}

func serverModeFromEnv(env string) string {
	serverMode := ServerModeDevelopment
	if env == "production" {
		serverMode = ServerModeProduction
	} else if env == "staging" {
		serverMode = ServerModeStaging
	}
	return serverMode
}

func (s Server) isProduction() bool {
	return s.mode == ServerModeProduction
}

func (s Server) isDevelopment() bool {
	return s.mode == ServerModeDevelopment
}

func (s Server) isStaging() bool {
	return s.mode == ServerModeStaging
}
