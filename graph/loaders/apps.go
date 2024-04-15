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

// appReader
type appReader struct {
}

// getApps implements a batch function that can retrieve many apps by ID,
// for use in a dataloader
func (u *appReader) getApps(ctx context.Context, appIDs []string) ([]*model.App, []error) {
	apps := make([]*model.App, len(appIDs))
	errs := make([]error, len(appIDs))

	ll, err := logger.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get logger: %w", err))
		return apps, errs
	}
	ll = ll.WithField("reader", "app").WithField("method", "getApps").WithField("appIDs", appIDs)
	ll.Debug("debug")

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

	if len(appIDs) == 1 {
		clientApp, _, err := doClient.Apps.Get(ctx, appIDs[0])
		if err != nil {
			return nil, []error{fmt.Errorf("unable to get app: %w", err)}
		}
		return []*model.App{model_helpers.AppFromGodo(clientApp)}, errs
	} else {
		// create options. initially, these will be blank
		opts := &godo.ListOptions{}
		for {
			ll.WithField("opts", opts).Info("doClient.Apps.List")
			// FIXME - call Get if its only a single one instead of list
			clientApps, resp, err := doClient.Apps.List(ctx, opts)
			if err != nil {
				return nil, []error{fmt.Errorf("unable to get apps: %w", err)}
			}

			for _, app := range clientApps {
				if pos, ok := appIDMap[app.ID]; ok {
					delete(appIDMap, app.ID)
					apps[pos] = model_helpers.AppFromGodo(app)
				}
			}

			if len(appIDMap) == 0 {
				// we got them all
				break
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

	for id, pos := range appIDMap {
		errs[pos] = fmt.Errorf("app %s is not found", id)
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
