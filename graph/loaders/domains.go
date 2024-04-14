package loaders

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"

	"github.com/halkeye/digitalocean-graphql-api/graph/digitalocean"
	"github.com/halkeye/digitalocean-graphql-api/graph/logger"
	"github.com/halkeye/digitalocean-graphql-api/graph/model"
)

// domainReader
type domainReader struct {
}

// getDomains implements a batch function that can retrieve many domains by ID,
// for use in a dataloader
func (u *domainReader) getDomains(ctx context.Context, domainIDs []string) ([]*model.Domain, []error) {
	domains := make([]*model.Domain, len(domainIDs))
	errs := make([]error, len(domainIDs))

	ll, err := logger.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get logger: %w", err))
	}
	ll = ll.WithField("reader", "domain").WithField("method", "getDomains").WithField("domainsIDs", domainIDs)
	ll.Info("debug")

	// stmt, err := u.db.PrepareContext(ctx, `SELECT id, name FROM domains WHERE id IN (?`+strings.Repeat(",?", len(domainIDs)-1)+`)`)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer stmt.Close()

	// rows, err := stmt.QueryContext(ctx, domainIDs)
	// if err != nil {
	// 	return nil, []error{err}
	// }
	// defer rows.Close()

	doClient, err := digitalocean.For(ctx)
	if err != nil {
		return nil, []error{fmt.Errorf("unable to get do client: %w", err)}
	}

	domainIDMap := map[string]int{}
	for pos, domainID := range domainIDs {
		domainIDMap[domainID] = pos
	}

	// create options. initially, these will be blank
	opts := &godo.ListOptions{}
	for {
		ll.WithField("opts", opts).Info("doClient.Domains.List")
		clientDomains, resp, err := doClient.Domains.List(ctx, opts)
		if err != nil {
			return nil, []error{fmt.Errorf("unable to get domains: %w", err)}
		}

		for _, domain := range clientDomains {
			if pos, ok := domainIDMap[domain.Name]; ok {
				delete(domainIDMap, domain.Name)
				domains[pos] = &model.Domain{
					ID:       domain.URN(),
					Name:     domain.Name,
					TTL:      domain.TTL,
					ZoneFile: &domain.ZoneFile,
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

	for id, pos := range domainIDMap {
		errs[pos] = fmt.Errorf("%s is not found", id)
	}

	return domains, errs
}

// GetDomain returns single domain by id efficiently
func GetDomain(ctx context.Context, domainID string) (*model.Domain, error) {
	loaders := For(ctx)
	return loaders.DomainLoader.Load(ctx, domainID)
}

// GetDomains returns many domains by ids efficiently
func GetDomains(ctx context.Context, domainIDs []string) ([]*model.Domain, error) {
	loaders := For(ctx)
	return loaders.DomainLoader.LoadAll(ctx, domainIDs)
}
