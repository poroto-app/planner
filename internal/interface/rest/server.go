package rest

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	port string
	mode string
}

const (
	ServerModeDevelopment = "development"
	ServerModeStaging     = "staging"
	ServerModeProduction  = "production"
)

func NewRestServer(env string) *Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{
		port: port,
		mode: serverModeFromEnv(env),
	}
}

func (s Server) ServeHTTP() error {
	r := gin.Default()

	if s.isProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST"},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowOriginFunc: func(origin string) bool {
			if !s.isProduction() {
				return true
			}
			protocol := os.Getenv("WEB_PROTOCOL")
			host := os.Getenv("WEB_HOST")
			return origin == fmt.Sprintf("%s://%s", protocol, host)
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from planner API",
		})
	})

	r.POST("/graphql", GraphQlQuery)
	if !s.isProduction() {
		r.GET("/graphql/playground", GraphQlPlayGround)
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
