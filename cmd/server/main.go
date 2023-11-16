package main

import (
	"log"
	"os"
	"poroto.app/poroto/planner/internal/env"

	"poroto.app/poroto/planner/internal/interface/rest"
)

func init() {
	env.LoadEnv()
}

func main() {
	s := rest.NewRestServer(os.Getenv("ENV"))
	if err := s.ServeHTTP(); err != nil {
		log.Fatalf("error while starting server: %v", err)
	}
}
