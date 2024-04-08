package loaders

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"

	"github.com/halkeye/digitalocean-graphql-api/graph/digitalocean"
	"github.com/halkeye/digitalocean-graphql-api/graph/model"
)

// domainReader
type domainReader struct {
}

// getDomains implements a batch function that can retrieve many domains by ID,
// for use in a dataloader
func (u *domainReader) getDomains(ctx context.Context, domainIDs []string) ([]*model.Domain, []error) {
	fmt.Printf("domainIDs: %v\n", domainIDs)

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
	domains := make([]*model.Domain, 0, len(domainIDs))
	errs := make([]error, 0, len(domainIDs))

	doClient, err := digitalocean.For(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to get do client: %w", err))
		return domains, errs
	}

	domainIDMap := map[string]int{}
	for pos, domainID := range domainIDs {
		domainIDMap[domainID] = pos
	}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		clientDomains, resp, err := doClient.Domains.List(ctx, opt)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to get do client: %w", err))
			return domains, errs
		}

		for _, domain := range clientDomains {
			if pos, ok := domainIDMap[domain.URN()]; ok {
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
			errs = append(errs, fmt.Errorf("unable to get current page: %w", err))
			return domains, errs

		}

		// set the page we want for the next request
		opt.Page = page + 1
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
