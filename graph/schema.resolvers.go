package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"strings"

	"github.com/halkeye/digitalocean-graphql-api/graph/loaders"
	"github.com/halkeye/digitalocean-graphql-api/graph/logger"
	"github.com/halkeye/digitalocean-graphql-api/graph/model"
)

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	ll, err := logger.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get do client: %w", err)
	}
	ll = ll.WithField("resolver", "query").WithField("id", id)
	ll.Info("debug")

	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("not a valid do urn")
	}

	if parts[0] != "do" {
		return nil, fmt.Errorf("not a valid do urn: namespace")
	}

	switch parts[1] {
	case "project":
		return loaders.GetProject(ctx, parts[2])
	case "droplet":
		return loaders.GetDroplet(ctx, parts[2])
	case "app":
		return loaders.GetApp(ctx, parts[2])
	case "volume":
		return nil, fmt.Errorf("projectResourceResolver.Resource - volume not implemented")
	case "domain":
		return loaders.GetDomain(ctx, parts[2])
	case "spaces":
		return nil, fmt.Errorf("projectResourceResolver.Resource - domain not implemented")
	default:
		return nil, fmt.Errorf("not a valid do urn: collection")
	}
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
