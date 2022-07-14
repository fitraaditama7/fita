package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"test-fite/graph/generated"
	"test-fite/graph/model"
)

// Checkout is the resolver for the checkout field.
func (r *mutationResolver) Checkout(ctx context.Context, input []*model.NewCheckout) (*model.ProductTotal, error) {
	return processCheckout(input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
