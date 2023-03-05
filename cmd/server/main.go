package main

import (
	"log"

	"poroto.app/poroto/planner/internal/interface/rest"
)

func main() {
	s := rest.NewRestServer(false)
	if err := s.ServeHTTP(); err != nil {
		log.Fatalf("error while starting server: %v", err)
	}
}
