package query

import (
	"fmt"
	"slices"

	"github.com/aquasecurity/esquery"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

func CreateQuery(filters *model.Filters, excludedKeys ...string) (*esquery.BoolQuery, error) {
	boolQuery := esquery.Bool()

	for _, filter := range filters.All() {
		if filter.IsEmpty() || slices.Contains(excludedKeys, filter.Key()) {
			continue
		}

		if builderByFilterKey := customBuilders[filter.Key()]; builderByFilterKey != nil {
			if ok := builderByFilterKey(filter, filters, boolQuery); !ok {
				return nil, fmt.Errorf("custom builder for field %s failed", filter.Key())
			}

			continue
		}

		opType := filter.OperationType()
		builderByOperationType := buildersByOperationType[opType]
		if builderByOperationType == nil {
			return nil, fmt.Errorf("builder for operation type %s not found", opType)
		}

		if ok := builderByOperationType(filter, filters, boolQuery); !ok {
			return nil, fmt.Errorf("builder for operation type %s and field %s failed", opType, filter.Key())
		}

	}

	return boolQuery, nil
}
