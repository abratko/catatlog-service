package service

import (
	"context"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
)

type OptFetcher interface {
	GetKey() string
	Exec(ctx context.Context, optDtos any) (any, error)
}

type filtersFactory struct {
	optsFetchersByType map[string]OptFetcher
	optsFetchersByKey  map[string]OptFetcher
}

func NewFiltersFactory(optsFetchers []OptFetcher) *filtersFactory {
	instance := &filtersFactory{
		optsFetchersByKey: map[string]OptFetcher{},
	}

	for _, fetcher := range optsFetchers {
		if fetcher.GetKey() == "" {
			panic("fetcher has no key")
		}

		instance.optsFetchersByKey[fetcher.GetKey()] = fetcher
	}

	return instance
}

func (f *filtersFactory) createFilter(filterDto dto.Filter) (*model.Filter, error) {
	opts := []func(*model.Filter) error{}
	if filterDto.Value != nil {
		opts = append(opts, model.WithValue(filterDto.Value))
	}

	if filterDto.DefaultValue != nil {
		opts = append(opts, model.WithDefaultValue(filterDto.DefaultValue))
	}

	filterOptsFetcher, ok := f.optsFetchersByType[filterDto.OperationType]
	if ok {
		opts = append(opts, model.WithFetchOpts(filterOptsFetcher.Exec))
	} else if filterKeyOptsFetcher, ok := f.optsFetchersByKey[filterDto.Key]; ok {
		opts = append(opts, model.WithFetchOpts(filterKeyOptsFetcher.Exec))
	}

	filter, err := model.NewFilter(filterDto, opts...)
	if err != nil {
		return nil, err
	}

	return filter, nil
}

func (f *filtersFactory) Exec(filterDtos []dto.Filter) (*model.Filters, error) {
	return model.NewFilters(filterDtos, f.createFilter)
}
