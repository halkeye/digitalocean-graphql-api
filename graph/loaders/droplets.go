package loaders

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"

	"github.com/halkeye/digitalocean-graphql-api/graph/digitalocean"
	"github.com/halkeye/digitalocean-graphql-api/graph/logger"
	"github.com/halkeye/digitalocean-graphql-api/graph/model"
)

// dropletReader
type dropletReader struct {
}

// getDroplets implements a batch function that can retrieve many droplets by ID,
// for use in a dataloader
func (u *dropletReader) getDroplets(ctx context.Context, dropletIDs []string) ([]*model.Droplet, []error) {
	// stmt, err := u.db.PrepareContext(ctx, `SELECT id, name FROM droplets WHERE id IN (?`+strings.Repeat(",?", len(dropletIDs)-1)+`)`)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer stmt.Close()

	// rows, err := stmt.QueryContext(ctx, dropletIDs)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer rows.Close()
	ll, err := logger.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get logger: %w", err)}
	}
	ll = ll.WithField("reader", "droplet").WithField("method", "getDroplets").WithField("dropletIDs", dropletIDs)
	ll.Info("debug")

	doClient, err := digitalocean.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get do client: %w", err)}
	}

	dropletIDMap := map[string]int{}
	for pos, dropletID := range dropletIDs {
		dropletIDMap[dropletID] = pos
	}

	droplets := make([]*model.Droplet, len(dropletIDs))
	errs := make([]error, len(dropletIDs))

	// create options. initially, these will be blank
	opts := &godo.ListOptions{}
	for {
		ll.WithField("opts", opts).Info("doClient.Droplets.List")
		clientDroplets, resp, err := doClient.Droplets.List(ctx, opts)
		if err != nil {
			return nil, []error{fmt.Errorf("unable to get droplets: %w", err)}
		}

		for _, droplet := range clientDroplets {
			if pos, ok := dropletIDMap[fmt.Sprint(droplet.ID)]; ok {
				delete(dropletIDMap, fmt.Sprint(droplet.ID))
				errs[pos] = nil
				droplets[pos] = &model.Droplet{
					ID:     droplet.URN(),
					Name:   droplet.Name,
					Memory: &droplet.Memory,
					Vcpus:  &droplet.Vcpus,
					Disk:   &droplet.Disk,
					Region: &model.Region{
						ID:   droplet.Region.Slug,
						Name: droplet.Region.Name,
					},
					SizeSlug:  &droplet.SizeSlug,
					BackupIDs: droplet.BackupIDs,
				}
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

	for id, pos := range dropletIDMap {
		errs[pos] = fmt.Errorf("%s is not found", id)
	}

	return droplets, errs
}

// GetDroplet returns single droplet by id efficiently
func GetDroplet(ctx context.Context, dropletID string) (*model.Droplet, error) {
	loaders := For(ctx)
	return loaders.DropletLoader.Load(ctx, dropletID)
}

// GetDroplets returns many droplets by ids efficiently
func GetDroplets(ctx context.Context, dropletIDs []string) ([]*model.Droplet, error) {
	loaders := For(ctx)
	return loaders.DropletLoader.LoadAll(ctx, dropletIDs)
}
