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

// dbaasReader
type dbaasReader struct {
}

// getDbaass implements a batch function that can retrieve many dbaass by ID,
// for use in a dataloader
func (u *dbaasReader) getDbaass(ctx context.Context, dbaasIDs []string) ([]*model.Dbaas, []error) {
	// stmt, err := u.db.PrepareContext(ctx, `SELECT id, name FROM dbaass WHERE id IN (?`+strings.Repeat(",?", len(dbaasIDs)-1)+`)`)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer stmt.Close()

	// rows, err := stmt.QueryContext(ctx, dbaasIDs)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer rows.Close()
	ll, err := logger.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get logger: %w", err)}
	}
	ll = ll.WithField("reader", "dbaas").WithField("method", "getDbaass").WithField("dbaasIDs", dbaasIDs)
	ll.Debug("debug")

	doClient, err := digitalocean.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get do client: %w", err)}
	}

	dbaasIDMap := map[string]int{}
	for pos, dbaasID := range dbaasIDs {
		dbaasIDMap[dbaasID] = pos
	}

	dbaass := make([]*model.Dbaas, len(dbaasIDs))
	errs := make([]error, len(dbaasIDs))

	if len(dbaasIDs) == 1 {
		ll.Info("looking up single database")
		dbaas, _, err := doClient.Databases.Get(ctx, dbaasIDs[0])
		if err != nil {
			return nil, []error{fmt.Errorf("unable to get dbaas: %w", err)}
		}
		return []*model.Dbaas{model_helpers.DbaasFromGodo(dbaas)}, errs
	} else {
		// create options. initially, these will be blank
		opts := &godo.ListOptions{
			PerPage: 100,
		}
		for {
			ll.WithField("opts", opts).Info("doClient.Dbaass.List")
			clientDbaass, resp, err := doClient.Databases.List(ctx, opts)
			if err != nil {
				return nil, []error{fmt.Errorf("unable to get dbaass: %w", err)}
			}

			for _, dbaas := range clientDbaass {
				ll.Info("looking at: " + dbaas.ID)
				if pos, ok := dbaasIDMap[fmt.Sprint(dbaas.ID)]; ok {
					delete(dbaasIDMap, fmt.Sprint(dbaas.ID))
					errs[pos] = nil
					dbaass[pos] = model_helpers.DbaasFromGodo(&dbaas)
				}
			}
			if len(dbaasIDs) == 0 {
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

	for id, pos := range dbaasIDMap {
		errs[pos] = fmt.Errorf("dbaas %s is not found", id)
		ll.WithError(errs[pos]).Info("dbaas not found")
	}

	return dbaass, errs
}

// GetDbaas returns single dbaas by id efficiently
func GetDbaas(ctx context.Context, dbaasID string) (*model.Dbaas, error) {
	loaders := For(ctx)
	return loaders.DbaasLoader.Load(ctx, dbaasID)
}

// GetDbaass returns many dbaass by ids efficiently
func GetDbaass(ctx context.Context, dbaasIDs []string) ([]*model.Dbaas, error) {
	loaders := For(ctx)
	return loaders.DbaasLoader.LoadAll(ctx, dbaasIDs)
}
