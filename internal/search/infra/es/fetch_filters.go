package es

import (
	"context"
	"fmt"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
)

var FilterConfigs = []dto.Filter{
	{
		Key:           "location",
		Label:         "Location",
		OperationType: model.TermsOp,
		Field:         "offers.facility_address.id",
	},
	{
		Key:           "category",
		Label:         "Category",
		OperationType: model.TermsOp,
		Field:         "category_ids",
	},
	{
		Key:           "condition",
		Label:         "Condition",
		OperationType: model.TermsOp,
		Field:         "condition",
	},
	{
		Key:           "manufacturer",
		Label:         "Manufacturer",
		OperationType: model.TermsOp,
		Field:         "brands_collection.title",
	},
	{
		Key:           "price",
		Label:         "Price",
		OperationType: model.RangeOp,
		Field:         "price",
	},
	{
		Key:           "inStock",
		Label:         "In Stock",
		OperationType: model.ToggleOp,
		Field:         "stock.qty",
		DefaultValue:  true,
	},
	{
		Key:           "text",
		Label:         "Search",
		OperationType: model.MatchOp,
		Field:         "name^5,description,sku",
	},
}

type createFilters = func(filterDtos []dto.Filter) (*model.Filters, error)
type fetchFilters struct {
	createFilters createFilters
}

func NewFetchFilters(
	createFilters createFilters,
) *fetchFilters {
	return &fetchFilters{
		createFilters: createFilters,
	}
}

func (f *fetchFilters) Exec(
	ctx context.Context,
	inputFilters map[string]any,
) (*model.Filters, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	var filtersDto []dto.Filter
	for _, dtoFilter := range FilterConfigs {
		filterValue := inputFilters[dtoFilter.Key]
		if filterValue != nil {
			dtoFilter.Value = filterValue
		}

		filtersDto = append(filtersDto, dtoFilter)
	}

	filters, err := f.createFilters(filtersDto)
	if err != nil {
		return nil, err
	}

	return filters, nil
}
