package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() {
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
