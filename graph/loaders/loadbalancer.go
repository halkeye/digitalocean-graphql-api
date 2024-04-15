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

// loadbalancerReader
type loadbalancerReader struct {
}

// getLoadBalancers implements a batch function that can retrieve many loadbalancers by ID,
// for use in a dataloader
func (u *loadbalancerReader) getLoadBalancers(ctx context.Context, loadbalancerIDs []string) ([]*model.LoadBalancer, []error) {
	// stmt, err := u.db.PrepareContext(ctx, `SELECT id, name FROM loadbalancers WHERE id IN (?`+strings.Repeat(",?", len(loadbalancerIDs)-1)+`)`)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer stmt.Close()

	// rows, err := stmt.QueryContext(ctx, loadbalancerIDs)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer rows.Close()
	ll, err := logger.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get logger: %w", err)}
	}
	ll = ll.WithField("reader", "loadbalancer").WithField("method", "getLoadBalancers").WithField("loadbalancerIDs", loadbalancerIDs)
	ll.Debug("debug")

	doClient, err := digitalocean.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get do client: %w", err)}
	}

	loadbalancerIDMap := map[string]int{}
	for pos, loadbalancerID := range loadbalancerIDs {
		loadbalancerIDMap[loadbalancerID] = pos
	}

	loadbalancers := make([]*model.LoadBalancer, len(loadbalancerIDs))
	errs := make([]error, len(loadbalancerIDs))

	if len(loadbalancerIDs) == 1 {
		loadbalancer, _, err := doClient.LoadBalancers.Get(ctx, loadbalancerIDs[0])
		if err != nil {
			return nil, []error{fmt.Errorf("unable to get loadbalancer: %w", err)}
		}
		return []*model.LoadBalancer{model_helpers.LoadBalancerFromGodo(loadbalancer)}, errs
	} else {
		// create options. initially, these will be blank
		opts := &godo.ListOptions{}
		for {
			ll.WithField("opts", opts).Info("doClient.LoadBalancers.List")
			clientLoadBalancers, resp, err := doClient.LoadBalancers.List(ctx, opts)
			if err != nil {
				return nil, []error{fmt.Errorf("unable to get loadbalancers: %w", err)}
			}

			for _, loadbalancer := range clientLoadBalancers {
				if pos, ok := loadbalancerIDMap[fmt.Sprint(loadbalancer.ID)]; ok {
					delete(loadbalancerIDMap, fmt.Sprint(loadbalancer.ID))
					errs[pos] = nil
					loadbalancers[pos] = model_helpers.LoadBalancerFromGodo(&loadbalancer)
				}
			}
			if len(loadbalancerIDMap) == 0 {
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

	for id, pos := range loadbalancerIDMap {
		errs[pos] = fmt.Errorf("loadbalancer %s is not found", id)
	}

	return loadbalancers, errs
}

// GetLoadBalancer returns single loadbalancer by id efficiently
func GetLoadBalancer(ctx context.Context, loadbalancerID string) (*model.LoadBalancer, error) {
	loaders := For(ctx)
	return loaders.LoadBalancerLoader.Load(ctx, loadbalancerID)
}

// GetLoadBalancers returns many loadbalancers by ids efficiently
func GetLoadBalancers(ctx context.Context, loadbalancerIDs []string) ([]*model.LoadBalancer, error) {
	loaders := For(ctx)
	return loaders.LoadBalancerLoader.LoadAll(ctx, loadbalancerIDs)
}
