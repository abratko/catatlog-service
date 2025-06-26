package app

import (
	"context"
	"fmt"
	"sync"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/errs"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

type fetchFilters = func(context.Context, map[string]any) (*model.Filters, error)
type fetchFilterOpts = func(context.Context, *model.Filters) (map[string]any, error)

type FetchFilterWithOpts struct {
	fetchFilters    fetchFilters
	fetchFilterOpts fetchFilterOpts
}

func NewFetchFiltersWithOpts(
	fetchFilters fetchFilters,
	fetchFilterOpts fetchFilterOpts,
) *FetchFilterWithOpts {
	return &FetchFilterWithOpts{
		fetchFilters:    fetchFilters,
		fetchFilterOpts: fetchFilterOpts,
	}
}

func (f *FetchFilterWithOpts) Exec(
	ctx context.Context,
	inputFilterVals map[string]any,
) (*model.Filters, error) {

	filters, err := f.fetchFilters(ctx, inputFilterVals)
	if err != nil {
		return nil, err
	}

	filterOptDtos, err := f.fetchFilterOpts(ctx, filters)
	if err != nil {
		return nil, err
	}

	err = f.initFilterOpts(ctx, filters, filterOptDtos)
	if err != nil {
		return nil, err
	}

	return filters, nil
}

func (f *FetchFilterWithOpts) initFilterOpts(
	ctx context.Context,
	filters *model.Filters,
	optDtos map[string]any,
) error {
	var wg sync.WaitGroup
	errsCh := make(chan error)
	for _, filter := range filters.OnlyFaceted() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := filter.InitOpts(ctx, optDtos[filter.Key()]); err != nil {
				errsCh <- fmt.Errorf("%w: %s", errs.ErrInternal, err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errsCh)
	}()

	err := <-errsCh
	return err
}
