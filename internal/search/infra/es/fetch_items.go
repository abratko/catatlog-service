package es

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/request"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/result"
)

type execRequest = func(context.Context, map[string]any) (*esapi.Response, error)

type FetchItems struct {
	execRequest execRequest
}

func NewFetchItems(
	execRequest execRequest,
) *FetchItems {
	return &FetchItems{
		execRequest: execRequest,
	}
}

func (f *FetchItems) Exec(
	ctx context.Context,
	criteria *model.SearchCriteria,
) (*dto.ItemsCollection, error) {
	requestDto, err := request.CreateFetchItems(criteria)
	if err != nil {
		return nil, err
	}

	res, err := f.execRequest(ctx, requestDto)
	if err != nil {
		return nil, err
	}

	result, err := result.Create(res)
	if err != nil {
		return nil, err
	}

	return result, nil
}
