package agg

import "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"

type M = map[string]any

func Location(
	filter *model.Filter,
	filters *model.Filters,
) (map[string]any, error) {

	minStock := 0
	if filters.Get("inStock").Value().(bool) {
		minStock = 1
	}

	condition := M{
		"nested": M{
			"path": "offers",
		},
		"aggs": M{
			filter.Key(): M{
				"filter": M{
					"bool": M{
						"filter": M{
							"range": M{
								"offers.stock.qty": M{
									"gte": minStock,
								},
							},
						},
					},
				},
				"aggs": M{
					filter.Key(): M{
						"terms": M{
							"field": filter.Field(),
						},
					},
				},
			},
		},
	}

	return condition, nil
}
