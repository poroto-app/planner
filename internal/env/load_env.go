package env

import (
	"fmt"
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

	if err := godotenv.Load(fmt.Sprintf(".env.%s", env)); err != nil {
		log.Fatalf("error while loading .env.%s: %v", env, err)
	}

	if err := godotenv.Load(fmt.Sprintf(".env.%s.local", env)); err != nil {
		log.Fatalf("error while loading .env: %v", err)
	}
}
