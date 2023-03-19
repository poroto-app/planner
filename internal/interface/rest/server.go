package rest

import (
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port       string
	production bool
}

func NewRestServer(production bool) *Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{
		port:       port,
		production: production,
	}
}

func (s Server) ServeHTTP() error {
	r := gin.Default()

	if s.production {
		gin.SetMode(gin.ReleaseMode)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from planner API",
		})
	})

	if err := r.Run(":" + s.port); err != nil {
		return nil
	}

	return nil
}
