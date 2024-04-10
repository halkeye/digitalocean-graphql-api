package loaders

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/google/uuid"

	"github.com/halkeye/digitalocean-graphql-api/graph/digitalocean"
	"github.com/halkeye/digitalocean-graphql-api/graph/logger"
	"github.com/halkeye/digitalocean-graphql-api/graph/model"
)

// appReader
type appReader struct {
}

// getApps implements a batch function that can retrieve many apps by ID,
// for use in a dataloader
func (u *appReader) getApps(ctx context.Context, appIDs []string) ([]*model.App, []error) {
	apps := make([]*model.App, 0, len(appIDs))
	errs := make([]error, 0, len(appIDs))

	ll, err := logger.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get do client: %w", err))
		return apps, errs
	}
	ll = ll.WithField("reader", "app").WithField("method", "getApps").WithField("appIDs", appIDs)
	ll.Info("debug")

	// stmt, err := u.db.PrepareContext(ctx, `SELECT id, name FROM apps WHERE id IN (?`+strings.Repeat(",?", len(appIDs)-1)+`)`)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer stmt.Close()

	// rows, err := stmt.QueryContext(ctx, appIDs)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer rows.Close()
	doClient, err := digitalocean.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get do client: %w", err))
		return apps, errs
	}

	appIDMap := map[string]int{}
	for pos, appID := range appIDs {
		appIDMap[appID] = pos
	}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		ll.WithField("opt", opt).Info("doClient.Apps.List")
		clientApps, resp, err := doClient.Apps.List(ctx, opt)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to get do client: %w", err))
			return apps, errs
		}

		for _, app := range clientApps {
			if pos, ok := appIDMap[app.URN()]; ok {
				apps[pos] = &model.App{
					ID:                     app.URN(),
					Owner:                  &model.Team{UUID: uuid.MustParse(app.OwnerUUID)},
					LastDeploymentActiveAt: &app.LastDeploymentActiveAt,
					DefaultIngress:         &app.DefaultIngress,
					CreatedAt:              &app.CreatedAt,
					UpdatedAt:              &app.UpdatedAt,
				}
			}
		}
		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to get current page: %w", err))
			return apps, errs

		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return apps, errs
}

// GetApp returns single app by id efficiently
func GetApp(ctx context.Context, appID string) (*model.App, error) {
	loaders := For(ctx)
	return loaders.AppLoader.Load(ctx, appID)
}

// GetApps returns many apps by ids efficiently
func GetApps(ctx context.Context, appIDs []string) ([]*model.App, error) {
	loaders := For(ctx)
	return loaders.AppLoader.LoadAll(ctx, appIDs)
}
