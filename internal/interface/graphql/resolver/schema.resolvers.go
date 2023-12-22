package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"

	"poroto.app/poroto/planner/internal/interface/graphql/generated"
)

// Ping is the resolver for the ping field.
func (r *mutationResolver) Ping(ctx context.Context, message string) (string, error) {
	return message, nil
}

// Version is the resolver for the version field.
func (r *queryResolver) Version(ctx context.Context) (string, error) {
	return "0.0.1", nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }