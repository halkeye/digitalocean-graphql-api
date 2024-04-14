package loaders

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"

	"github.com/halkeye/digitalocean-graphql-api/graph/digitalocean"
	"github.com/halkeye/digitalocean-graphql-api/graph/logger"
	"github.com/halkeye/digitalocean-graphql-api/graph/model"
	"github.com/halkeye/digitalocean-graphql-api/graph/model_helpers"
)

// projectReader
type projectReader struct {
}

// getProjects implements a batch function that can retrieve many projects by ID,
// for use in a dataloader
func (u *projectReader) getProjects(ctx context.Context, projectIDs []string) ([]*model.Project, []error) {
	projects := make([]*model.Project, len(projectIDs))
	errs := make([]error, len(projectIDs))

	ll, err := logger.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get logger: %w", err))
		return projects, errs
	}
	ll = ll.WithField("reader", "project").WithField("method", "getProjects").WithField("projectIDs", projectIDs)
	ll.Info("debug")

	// stmt, err := u.db.PrepareContext(ctx, `SELECT id, name FROM projects WHERE id IN (?`+strings.Repeat(",?", len(projectIDs)-1)+`)`)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer stmt.Close()

	// rows, err := stmt.QueryContext(ctx, projectIDs)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer rows.Close()
	doClient, err := digitalocean.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get do client: %w", err))
		return projects, errs
	}

	projectIDMap := map[string]int{}
	for pos, projectID := range projectIDs {
		projectIDMap[projectID] = pos
	}

	if len(projectIDs) == 1 {
		clientProject, _, err := doClient.Projects.Get(ctx, projectIDs[0])
		if err != nil {
			return nil, []error{fmt.Errorf("unable to get projects: %w", err)}
		}
		return []*model.Project{model_helpers.ProjectFromGodo(clientProject)}, errs
	} else {
		// create options. initially, these will be blank
		opts := &godo.ListOptions{}
		for {
			ll.WithField("opts", opts).Info("doClient.Projects.List")
			// FIXME - call Get if its only a single one instead of list
			clientProjects, resp, err := doClient.Projects.List(ctx, opts)
			if err != nil {
				return nil, []error{fmt.Errorf("unable to get projects: %w", err)}
			}

			for _, project := range clientProjects {
				if pos, ok := projectIDMap[project.ID]; ok {
					delete(projectIDMap, project.ID)
					projects[pos] = model_helpers.ProjectFromGodo(&project)
				}
			}
			// if we are at the last page, break out the for loop
			if resp.Links == nil || resp.Links.IsLastPage() {
				break
			}

			page, err := resp.Links.CurrentPage()
			if err != nil {
				return nil, []error{fmt.Errorf("unable to get current page: %w", err)}
			}

			// set the page we want for the next request
			opts.Page = page + 1
		}
	}

	for id, pos := range projectIDMap {
		errs[pos] = fmt.Errorf("%s is not found", id)
	}

	return projects, errs
}

// GetProject returns single project by id efficiently
func GetProject(ctx context.Context, projectID string) (*model.Project, error) {
	loaders := For(ctx)
	return loaders.ProjectLoader.Load(ctx, projectID)
}

// GetProjects returns many projects by ids efficiently
func GetProjects(ctx context.Context, projectIDs []string) ([]*model.Project, error) {
	loaders := For(ctx)
	return loaders.ProjectLoader.LoadAll(ctx, projectIDs)
}
