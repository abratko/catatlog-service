package agg

import "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"

func Terms(
	filter *model.Filter,
	filters *model.Filters,
) (map[string]any, error) {

	agg := map[string]any{
		"terms": map[string]any{
			"field": filter.Field(),
			"size":  1000,
		},
	}

	return agg, nil
}
