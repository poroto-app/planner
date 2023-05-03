package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"poroto.app/poroto/planner/internal/interface/rest"
)

func init() {
	env := os.Getenv("ENV")
	if "" == env {
		env = "development"
	}

	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalf("error while loading .env.local: %v", err)
	}

	if err := godotenv.Load(".env." + env); err != nil {
		log.Fatalf("error while loading .env.%s: %v", env, err)
	}
}

func main() {
	s := rest.NewRestServer(os.Getenv("ENV"))
	if err := s.ServeHTTP(); err != nil {
		log.Fatalf("error while starting server: %v", err)
	}
}
