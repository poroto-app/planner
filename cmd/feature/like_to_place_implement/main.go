package main

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

func init() {
	env.LoadEnv()
}

func main() {

	ctx := context.Background()
	rep, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		fmt.Printf("Error creating plan candidate repository: %v\n", err)
	}

	planCandidateId := "291de5a4-f96a-4061-bb3b-c139d56f6ae0"
	placeId := "R4rEWuWpEEUmPtAZeh9H"

	err = rep.UpdateLikeToPlaceInPlanCandidate(ctx, planCandidateId, placeId)
	if err != nil {
		fmt.Printf("Error updating like to place in plan candidate: %v\n", err)
		return
	}

	fmt.Println("Like to place updated successfully!")
}
