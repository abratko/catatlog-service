package builders

import (
	"github.com/aquasecurity/esquery"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

func InStock(
	filter *model.Filter,
	filters *model.Filters,
	query *esquery.BoolQuery,
) bool {
	rangeQuery := esquery.Range(filter.Field())

	if filter.Value().(bool) {
		rangeQuery.Gte(1)
	} else {
		rangeQuery.Gte(0)
	}

	query.Filter(rangeQuery)

	return true
}
