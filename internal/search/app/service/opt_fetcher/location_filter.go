package opt_fetcher

import (
	"context"
	"errors"
	"strconv"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

type LocationRepo interface {
	FindByIds(ctx context.Context, ids []string) ([]dto.Location, error)
}

type locationOptFetcher struct {
	findByIds FindByIds[string, dto.Location]
}

// We don't use interface instead dto.Location because opt-fetcher is exclusive adapter for
// LocationRepo.FindByIds method. So, if we would need repository , we would need to create new
// adapter.
func NewLocationOptFetcher(findByIds FindByIds[string, dto.Location]) *locationOptFetcher {
	return &locationOptFetcher{
		findByIds: findByIds,
	}
}

func (f *locationOptFetcher) GetKey() string {
	return "location"
}

func (f *locationOptFetcher) Exec(ctx context.Context, originOptDtos any) (any, error) {
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
		newOpts[i] = optsByValue[strconv.Itoa(item.Id)]
		newOpts[i].Props = map[string]any{
			"name":    item.Name,
			"country": item.Country.Code,
		}
	}

	return newOpts, nil
}
