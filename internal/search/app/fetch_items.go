package app

import (
	"context"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
)

type fetchItems = func(context.Context, *model.SearchCriteria) (*dto.ItemsCollection, error)
type FetchItems struct {
	fetchFilters fetchFilters
	fetchItems   fetchItems
}

func NewFetchItems(
	fetchFilters fetchFilters,
	fetchItems fetchItems,
) *FetchItems {
	return &FetchItems{
		fetchFilters: fetchFilters,
		fetchItems:   fetchItems,
	}
}

func (f *FetchItems) Exec(
	ctx context.Context,
	dtoCriteria dto.SearchCriteria,
) (*dto.ItemsCollection, error) {
	filters, err := f.fetchFilters(ctx, dtoCriteria.Filters)
	if err != nil {
		return nil, err
	}

	criteria := model.NewSearchCriteria(
		model.WithFilters(filters),
		model.WithOrderBy(dtoCriteria.OrderBy),
		model.WithPagination(dtoCriteria.Pagination),
	)

	items, err := f.fetchItems(ctx, criteria)
	if err != nil {
		return nil, err
	}

	return items, nil
}
