package agg

import "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"

func Range(
	filter *model.Filter,
	filters *model.Filters,
) (map[string]any, error) {
	agg := map[string]any{
		"stats": map[string]any{
			"field": filter.Field(),
		},
	}

	return agg, nil
}
