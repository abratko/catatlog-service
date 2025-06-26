package opt_fetcher

import (
	"context"
	"errors"

	category "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

type FindByIds[Id any, Res any] func(ctx context.Context, ids []Id) ([]Res, error)

type categoryOptFetcher struct {
	findByIds FindByIds[string, category.Category]
}

func NewCategoryOptFetcher(findByIds FindByIds[string, category.Category]) *categoryOptFetcher {
	return &categoryOptFetcher{
		findByIds: findByIds,
	}
}

func (f *categoryOptFetcher) GetKey() string {
	return "category"
}

func (f *categoryOptFetcher) Exec(ctx context.Context, originOptDtos any) (any, error) {
	optDtos, ok := originOptDtos.([]model.TermOpt)
	if !ok {
		return nil, errors.New("invalid optDtos")
	}

	optsByValue := map[string]model.TermOpt{}
	ids := make([]string, 0, len(optDtos))
	for _, opt := range optDtos {
		ids = append(ids, opt.Value)
		optsByValue[opt.Value] = opt
	}

	items, err := f.findByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	newOpts := make([]model.TermOpt, len(items))
	for i, item := range items {
		newOpts[i] = optsByValue[item.Id]
		newOpts[i].Props = map[string]any{
			"name": item.Name,
			"path": item.Path,
			"slug": item.Slug,
		}
	}

	return newOpts, nil
}
