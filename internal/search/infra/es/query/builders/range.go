package builders

import (
	"github.com/aquasecurity/esquery"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

func Range(
	filter *model.Filter,
	filters *model.Filters,
	query *esquery.BoolQuery,
) bool {
	rangeVal, ok := filter.Value().(model.RangeVal)
	if !ok {
		return false
	}

	from := rangeVal.From
	to := rangeVal.To

	if from == nil && to == nil {
		return false
	}

	if from != nil && to != nil && from == to {
		return false
	}

	rangeQuery := esquery.Range(filter.Field())

	if from != nil {
		rangeQuery.Gte(from)
	}

	if to != nil {
		rangeQuery.Lte(to)
	}

	query.Filter(rangeQuery)

	return true
}
