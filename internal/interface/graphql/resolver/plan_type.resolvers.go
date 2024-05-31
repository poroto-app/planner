package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/interface/graphql/generated"
	"poroto.app/poroto/planner/internal/interface/graphql/model"
)

// Collage is the resolver for the collage field.
func (r *planResolver) Collage(ctx context.Context, obj *model.Plan) (*model.PlanCollage, error) {
	panic(fmt.Errorf("not implemented: Collage - collage"))
}

// Plan returns generated.PlanResolver implementation.
func (r *Resolver) Plan() generated.PlanResolver { return &planResolver{r} }

type planResolver struct{ *Resolver }
