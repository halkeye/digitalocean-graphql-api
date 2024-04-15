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

// volumeReader
type volumeReader struct {
}

// getVolumes implements a batch function that can retrieve many volumes by ID,
// for use in a dataloader
func (u *volumeReader) getVolumes(ctx context.Context, volumeIDs []string) ([]*model.Volume, []error) {
	volumes := make([]*model.Volume, len(volumeIDs))
	errs := make([]error, len(volumeIDs))

	ll, err := logger.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get logger: %w", err))
		return volumes, errs
	}
	ll = ll.WithField("reader", "app").WithField("method", "getVolumes").WithField("volumeIDs", volumeIDs)
	ll.Debug("debug")

	// stmt, err := u.db.PrepareContext(ctx, `SELECT id, name FROM volumes WHERE id IN (?`+strings.Repeat(",?", len(volumeIDs)-1)+`)`)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer stmt.Close()

	// rows, err := stmt.QueryContext(ctx, volumeIDs)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer rows.Close()
	doClient, err := digitalocean.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get do client: %w", err))
		return volumes, errs
	}

	appIDMap := map[string]int{}
	for pos, appID := range volumeIDs {
		appIDMap[appID] = pos
	}

	if len(volumeIDs) == 1 {
		clientVolume, _, err := doClient.Storage.GetVolume(ctx, volumeIDs[0])
		if err != nil {
			return nil, []error{fmt.Errorf("unable to get volumes: %w", err)}
		}
		return []*model.Volume{model_helpers.VolumeFromGodo(clientVolume)}, errs
	} else {
		// create options. initially, these will be blank
		opts := &godo.ListVolumeParams{}
		for {
			ll.WithField("opts", opts).Info("doClient.Storage.ListVolumes")
			// FIXME - call Get if its only a single one instead of list
			clientVolumes, resp, err := doClient.Storage.ListVolumes(ctx, opts)
			if err != nil {
				return nil, []error{fmt.Errorf("unable to get volumes: %w", err)}
			}

			for _, volume := range clientVolumes {
				if pos, ok := appIDMap[volume.ID]; ok {
					delete(appIDMap, volume.ID)
					volumes[pos] = model_helpers.VolumeFromGodo(&volume)
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
			opts.ListOptions.Page = page + 1
		}
	}

	for id, pos := range appIDMap {
		errs[pos] = fmt.Errorf("volume %s is not found", id)
	}

	return volumes, errs
}

// GetVolume returns single app by id efficiently
func GetVolume(ctx context.Context, appID string) (*model.Volume, error) {
	loaders := For(ctx)
	return loaders.VolumeLoader.Load(ctx, appID)
}

// GetVolumes returns many volumes by ids efficiently
func GetVolumes(ctx context.Context, volumeIDs []string) ([]*model.Volume, error) {
	loaders := For(ctx)
	return loaders.VolumeLoader.LoadAll(ctx, volumeIDs)
}
