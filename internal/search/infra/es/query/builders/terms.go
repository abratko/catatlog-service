package builders

import (
	"github.com/aquasecurity/esquery"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

func Terms(
	filter *model.Filter,
	filters *model.Filters,
	query *esquery.BoolQuery,
) bool {
	values := filter.Value()
	if values == nil {
		return false
	}

	anySlice := convertToAnySlice(values)
	if len(anySlice) == 0 {
		return false
	}

	query.Filter(esquery.Terms(filter.Field(), anySlice...))

	return true
}

func convertToAnySlice(val any) []any {
	switch v := val.(type) {
	case []int:
		anySlice := make([]any, len(v))
		for i := 0; i < len(v); i++ {
			anySlice[i] = v[i]
		}

		return anySlice

	case []float64:
		anySlice := make([]any, len(v))
		for i := 0; i < len(v); i++ {
			anySlice[i] = v[i]
		}

		return anySlice

	case []string:
		anySlice := make([]any, len(v))
		for i := 0; i < len(v); i++ {
			anySlice[i] = v[i]
		}

		return anySlice

	default:
		return []any{}
	}
}
