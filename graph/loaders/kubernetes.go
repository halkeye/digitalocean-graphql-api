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

// kubernetesReader
type kubernetesReader struct {
}

// getKubernetesClusters implements a batch function that can retrieve many kubernetess by ID,
// for use in a dataloader
func (u *kubernetesReader) getKubernetesClusters(ctx context.Context, kubernetesClusterIDs []string) ([]*model.KubernetesCluster, []error) {
	ll, err := logger.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get logger: %w", err)}
	}
	ll = ll.WithField("reader", "kubernetes").WithField("method", "getKubernetesClusters").WithField("kubernetesClusterIDs", kubernetesClusterIDs)
	ll.Info("debug")

	doClient, err := digitalocean.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get do client: %w", err)}
	}

	kubernetesIDMap := map[string]int{}
	for pos, kubernetesID := range kubernetesClusterIDs {
		kubernetesIDMap[kubernetesID] = pos
	}

	kubernetess := make([]*model.KubernetesCluster, len(kubernetesClusterIDs))
	errs := make([]error, len(kubernetesClusterIDs))

	if len(kubernetesClusterIDs) == 1 {
		ll.Info("looking up single database")
		kubernetesCluster, _, err := doClient.Kubernetes.Get(ctx, kubernetesClusterIDs[0])
		if err != nil {
			return nil, []error{fmt.Errorf("unable to get kubernetes: %w", err)}
		}
		return []*model.KubernetesCluster{model_helpers.KubernetesClusterFromGodo(kubernetesCluster)}, errs
	} else {
		// create options. initially, these will be blank
		opts := &godo.ListOptions{}
		for {
			ll.WithField("opts", opts).Info("doClient.Kubernetes.List")
			clientClusters, resp, err := doClient.Kubernetes.List(ctx, opts)
			if err != nil {
				return nil, []error{fmt.Errorf("unable to get kubernetesClusters: %w", err)}
			}

			for _, kubernetes := range clientClusters {
				ll.Info("looking at: " + kubernetes.ID)
				if pos, ok := kubernetesIDMap[fmt.Sprint(kubernetes.ID)]; ok {
					delete(kubernetesIDMap, fmt.Sprint(kubernetes.ID))
					errs[pos] = nil
					kubernetess[pos] = model_helpers.KubernetesClusterFromGodo(kubernetes)
				}
			}

			if len(kubernetesIDMap) == 0 {
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

	for id, pos := range kubernetesIDMap {
		errs[pos] = fmt.Errorf("kubernetes %s is not found", id)
	}

	return kubernetess, errs
}

// GetKubernetesCluster returns single kubernetes by id efficiently
func GetKubernetesCluster(ctx context.Context, kubernetesID string) (*model.KubernetesCluster, error) {
	loaders := For(ctx)
	return loaders.KubernetesLoader.Load(ctx, kubernetesID)
}

// GetKubernetesClusters returns many kubernetess by ids efficiently
func GetKubernetesClusters(ctx context.Context, kubernetesIDs []string) ([]*model.KubernetesCluster, error) {
	loaders := For(ctx)
	return loaders.KubernetesLoader.LoadAll(ctx, kubernetesIDs)
}
