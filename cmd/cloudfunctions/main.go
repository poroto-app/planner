package main

import (
	"log"
	"os"
	"poroto.app/poroto/planner/internal/env"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"poroto.app/poroto/planner/internal/interface/cloudfunctions"
)

func init() {
	env.LoadEnv()

	functions.HTTP("DeleteExpiredPlanCandidates", cloudfunctions.DeleteExpiredPlanCandidates)
}

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
