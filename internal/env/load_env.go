package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type LoadEnvOptions struct {
	SkipLoadErrors bool
}

type LoadEnvOption func(*LoadEnvOptions)

func WithSkipErrors() LoadEnvOption {
	return func(o *LoadEnvOptions) {
		o.SkipLoadErrors = true
	}
}

func LoadEnv(opts ...LoadEnvOption) {
	options := &LoadEnvOptions{}
	for _, opt := range opts {
		opt(options)
	}

	env := os.Getenv("ENV")
	if "" == env {
		env = "development"
		os.Setenv("ENV", env)
	}

	loadEnvFile := func(filename string) {
		if err := godotenv.Load(filename); err != nil {
			if options.SkipLoadErrors {
				log.Printf("error while loading %s: %v", filename, err)
			} else {
				log.Fatalf("error while loading %s: %v", filename, err)
			}
		}
	}

	loadEnvFile(".env.local")
	loadEnvFile(fmt.Sprintf(".env.%s.local", env))
	loadEnvFile(fmt.Sprintf(".env.%s", env))
}
