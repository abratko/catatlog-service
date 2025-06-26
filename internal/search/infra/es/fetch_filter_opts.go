package es

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/request"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/result"
)

type Fetch = func(ctx context.Context, queryDto map[string]any) (*esapi.Response, error)

type fetchFilterOpts struct {
	fetch Fetch
}

func NewFetchFilterOpts(
	fetch Fetch,
) *fetchFilterOpts {
	return &fetchFilterOpts{
		fetch: fetch,
	}
}

func (f *fetchFilterOpts) Exec(
	ctx context.Context,
	filters *model.Filters,
) (map[string]any, error) {
	requestDto, err := request.CreateFetchFilterOpts(filters)
	if err != nil {
		return nil, err
	}

	res, err := f.fetch(ctx, requestDto)
	if err != nil {
		return nil, err
	}

	result, err := result.CreateOptDtos(res)
	if err != nil {
		return nil, err
	}

	return result, nil
}
