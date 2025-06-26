package agg

import (
	"fmt"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

type AggFactoryMap map[string]func(filter *model.Filter, filters *model.Filters) (map[string]any, error)

var aggFactoriesByKey = AggFactoryMap{
	"location": Location,
}

func CreateAggs(filter *model.Filter, filters *model.Filters) (map[string]any, error) {
	if factory, ok := aggFactoriesByKey[filter.Key()]; ok {
		return factory(filter, filters)
	}

	if filter.OperationType() == model.TermsOp {
		return Terms(filter, filters)
	}

	if filter.OperationType() == model.RangeOp {
		return Range(filter, filters)
	}

	return nil, fmt.Errorf("not found aggregation factory by filter key: %s", filter.Key())
}
