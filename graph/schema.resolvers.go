package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/halkeye/digitalocean-graphql-api/graph/digitalocean"
	"github.com/halkeye/digitalocean-graphql-api/graph/loaders"
	"github.com/halkeye/digitalocean-graphql-api/graph/logger"
	"github.com/halkeye/digitalocean-graphql-api/graph/model"
	"github.com/halkeye/digitalocean-graphql-api/graph/model_helpers"
	"github.com/sirupsen/logrus"
)

// Resources is the resolver for the resources field.
func (r *projectResolver) Resources(ctx context.Context, obj *model.Project, first *int, after *string) (*model.ProjectResourcesConnection, error) {
	doClient, err := digitalocean.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get do client: %w", err)
	}

	ll, err := logger.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get logger: %w", err)
	}

	ll = ll.WithField("resolver", "Resources").WithField("parent.id", obj.ID).WithField("first", *first)
	ll.Info("debug")

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: *first,
	}

	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		opts.Page, err = strconv.Atoi(string(b))
		if err != nil {
			return nil, fmt.Errorf("unable to process cursor: %w", err)
		}
	}

	// FIXME - cap first
	edges := make([]*model.ProjectResourcesEdge, *first)
	count := 0

	ll.WithField("project.id", strings.Replace(obj.ID, "do:project:", "", 1)).WithField("opts", opts).Info("doClient.Projects.ListResources")

	projectResources, resp, err := doClient.Projects.ListResources(ctx, strings.Replace(obj.ID, "do:project:", "", 1), opts)
	if err != nil {
		return nil, fmt.Errorf("unable to get projects: %w", err)
	}

	ll.WithField("projectResources.length", len(projectResources)).Info("length")

	for _, pr := range projectResources {
		assignedAt, err := time.Parse(time.RFC3339, pr.AssignedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to parse assignedAt at: %w", err)
		}

		id := fmt.Sprintf(`do:projectresource:%s:%s`, obj.ID, pr.URN)
		edges[count] = &model.ProjectResourcesEdge{
			Cursor: base64.StdEncoding.EncodeToString([]byte(id)),
			Node: &model.ProjectResource{
				ID:         id,
				AssignedAt: assignedAt,
				Status:     pr.Status,
			},
		}
		count++
	}

	mc := &model.ProjectResourcesConnection{
		Edges: edges[:count],
		PageInfo: &model.PageInfo{
			HasPreviousPage: opts.Page != 1,
			HasNextPage:     !resp.Links.IsLastPage(),
		},
	}
	if mc.PageInfo.HasPreviousPage {
		mc.PageInfo.StartCursor = mustStringPtr(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", opts.Page))))
	}
	if mc.PageInfo.HasNextPage {
		mc.PageInfo.EndCursor = mustStringPtr(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", opts.Page+1))))
	}

	return mc, nil
}

// Resource is the resolver for the resource field.
func (r *projectResourceResolver) Resource(ctx context.Context, obj *model.ProjectResource) (model.Resource, error) {
	ll, err := logger.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get logger: %w", err)
	}
	ll = ll.WithField("resolver", "resolver").WithField("parent.id", obj.ID)
	ll.Debug("debug")

	urn := obj.ID
	if strings.HasPrefix(urn, "do:projectresource:") {
		// get rid of do:projectresource:do:project:uuid
		urn = strings.Join(strings.Split(urn, ":")[5:], ":")
	}

	return GetResource(ctx, ll, urn)
}

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	ll, err := logger.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get do client: %w", err)
	}
	ll = ll.WithField("resolver", "query").WithField("id", id)
	ll.Debug("debug")

	objtype, id, err := fromDoURN(id)
	if err != nil {
		return nil, err
	}

	if objtype == "project" {
		return loaders.GetProject(ctx, id)
	}

	return GetResource(ctx, ll, id)
}

// Projects is the resolver for the projects field.
func (r *queryResolver) Projects(ctx context.Context, first *int, after *string, last *int, before *string) (*model.ProjectsConnection, error) {
	// FIXME - ordering should be consistent -
	// https://relay.dev/graphql/connections.htm#sec-Edge-order
	// You may order the edges however your business logic dictates, and may determine the ordering based upon additional arguments not covered by this specification. But the ordering must be consistent from page to page, and importantly, The ordering of edges should be the same when using first/after as when using last/before, all other arguments being equal. It should not be reversed when using last/before. More formally:
	// When before: cursor is used, the edge closest to cursor must come last in the result edges.
	// When after: cursor is used, the edge closest to cursor must come first in the result edges.

	doClient, err := digitalocean.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get do client: %w", err)
	}

	ll, err := logger.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get logger: %w", err)
	}
	ll = ll.WithField("resolver", "projects")

	if first == nil {
		first = new(int)
		*first = 10
	}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: *first,
	}

	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		opts.Page, err = strconv.Atoi(string(b))
		if err != nil {
			return nil, fmt.Errorf("unable to process cursor: %w", err)
		}
	}
	ll = ll.WithFields(logrus.Fields{
		"first":        *first,
		"opts.page":    opts.Page,
		"opts.perpage": opts.PerPage,
	})
	ll.Debug("debug")

	edges := make([]*model.ProjectsEdge, *first)
	count := 0

	projects, resp, err := doClient.Projects.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("unable to get projects: %w", err)
	}

	for _, p := range projects {
		node := model_helpers.ProjectFromGodo(&p)
		edges[count] = &model.ProjectsEdge{
			Cursor: base64.StdEncoding.EncodeToString([]byte(node.ID)),
			Node:   node,
		}
		count++
	}

	page := 0
	if resp.Links == nil {
		page, err = resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("unable to get current page: %w", err)
		}
	}
	ll.WithField("page", page).Info("next page")

	mc := &model.ProjectsConnection{
		Edges: edges[:count],
		PageInfo: &model.PageInfo{
			HasPreviousPage: opts.Page != 1,
			HasNextPage:     !resp.Links.IsLastPage(),
		},
	}
	if mc.PageInfo.HasPreviousPage {
		mc.PageInfo.StartCursor = mustStringPtr(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", opts.Page-1))))
	}
	if mc.PageInfo.HasNextPage {
		mc.PageInfo.EndCursor = mustStringPtr(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", opts.Page+1))))
	}

	return mc, nil
}

// Account is the resolver for the account field.
func (r *queryResolver) Account(ctx context.Context) (*model.Account, error) {
	doClient, err := digitalocean.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get do client: %w", err)
	}

	ll, err := logger.For(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get logger: %w", err)
	}
	ll = ll.WithField("resolver", "account")
	ll.Info("get account")

	clientAccount, _, err := doClient.Account.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get accounts: %w", err)
	}
	return model_helpers.AccountFromGodo(clientAccount)
}

// Project returns ProjectResolver implementation.
func (r *Resolver) Project() ProjectResolver { return &projectResolver{r} }

// ProjectResource returns ProjectResourceResolver implementation.
func (r *Resolver) ProjectResource() ProjectResourceResolver { return &projectResourceResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type projectResolver struct{ *Resolver }
type projectResourceResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
