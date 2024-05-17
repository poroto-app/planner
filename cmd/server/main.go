package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"

	"poroto.app/poroto/planner/internal/interface/rest"
)

func init() {
	env.LoadEnv()
}

func main() {
	db, err := rdb.InitDB(false)
	if err != nil {
		log.Fatalf("error while initializing db: %v", err)
	}

	s, err := rest.NewRestServer(context.Background(), db, os.Getenv("ENV"))
	if err != nil {
		log.Fatalf("error while initializing server: %v", err)
	}

	if err := s.ServeHTTP(db); err != nil {
		log.Fatalf("error while starting server: %v", err)
	}
}
