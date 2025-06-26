package model

import "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"

type Filters struct {
	allAsMap          map[string]*Filter
	allAsSlice        []*Filter
	facetedFilters    []*Filter
	notFacetedFilters []*Filter
	facetedKeys       []string
	notFacetedKeys    []string
}

type createFilter = func(dto dto.Filter) (*Filter, error)

func NewFilters(
	dtos []dto.Filter,
	createFilter createFilter,
) (*Filters, error) {
	filters := make(map[string]*Filter)
	facetedFilters := make([]*Filter, 0)
	notFacetedFilters := make([]*Filter, 0)
	allAsSlice := make([]*Filter, 0)
	facetedKeys := make([]string, 0)
	notFacetedKeys := make([]string, 0)

	for _, dto := range dtos {
		filter, err := createFilter(dto)
		if err != nil {
			return nil, err
		}

		filters[filter.Key()] = filter
		allAsSlice = append(allAsSlice, filter)
		if filter.IsFaceted() {
			facetedFilters = append(facetedFilters, filter)
			facetedKeys = append(facetedKeys, filter.Key())

			continue
		}
		notFacetedFilters = append(notFacetedFilters, filter)
		notFacetedKeys = append(notFacetedKeys, filter.Key())

	}

	return &Filters{
		allAsMap:          filters,
		allAsSlice:        allAsSlice,
		facetedFilters:    facetedFilters,
		notFacetedFilters: notFacetedFilters,
		facetedKeys:       facetedKeys,
		notFacetedKeys:    notFacetedKeys,
	}, nil
}

func (f Filters) Get(key string) *Filter {
	return f.allAsMap[key]
}

func (f Filters) OnlyFaceted() []*Filter {
	return f.facetedFilters
}
func (f Filters) OnlyFacetedKeys() []string {
	return f.facetedKeys
}

func (f Filters) OnlyNotFacetedKeys() []string {
	return f.notFacetedKeys
}

func (f Filters) OnlyNotFaceted() []*Filter {
	return f.notFacetedFilters
}

func (f Filters) All() []*Filter {
	return f.allAsSlice
}
