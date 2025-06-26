package builders

import (
	"github.com/aquasecurity/esquery"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

type m = map[string]any

func Location(
	filter *model.Filter,
	filters *model.Filters,
	query *esquery.BoolQuery,
) bool {
	values := filter.Value()
	// the facility property inside product is nested document with self stock.qty property,
	// therefore we should select only facilities that satisfies inStock filter
	minStock := 0
	if filters.Get("inStock").Value().(bool) {
		minStock = 1
	}

	condition := m{
		"nested": m{
			"path": "offers",
			"query": m{
				"bool": m{
					"must": []m{
						{
							"terms": m{
								"offers.facility_address.id": values,
							},
						},
						{
							"range": m{
								"offers.stock.qty": m{
									"gte": minStock,
								},
							},
						},
					},
				},
			},
		},
	}

	query.Filter(esquery.CustomQuery(condition))

	return true
}
