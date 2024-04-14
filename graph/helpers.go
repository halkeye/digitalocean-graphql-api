package graph

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/halkeye/digitalocean-graphql-api/graph/loaders"
	"github.com/halkeye/digitalocean-graphql-api/graph/model"
)

func mustStringPtr(s string) *string {
	return &s
}

func fromDoURN(urn string) (string, string, error) {
	parts := strings.Split(urn, ":")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("not a valid do urn: %s", urn)
	}

	if parts[0] != "do" {
		return "", "", fmt.Errorf("urn has unhandled namespace: %s", urn)
	}

	return parts[1], parts[2], nil
}

func GetResource(ctx context.Context, ll *logrus.Entry, urn string) (model.Resource, error) {
	objtype, id, err := fromDoURN(urn)
	if err != nil {
		return nil, err
	}

	switch objtype {
	case "droplet":
		return loaders.GetDroplet(ctx, id)
	case "app":
		return loaders.GetApp(ctx, id)
	case "volume":
		return loaders.GetVolume(ctx, id)
	case "domain":
		return loaders.GetDomain(ctx, id)
	case "space":
		return nil, fmt.Errorf("projectResourceResolver.Resource - space not implemented")
	case "dbaas":
		return nil, fmt.Errorf("projectResourceResolver.Resource - dbaas not implemented")
	default:
		return nil, fmt.Errorf("no handler for %s", objtype)
	}
}
