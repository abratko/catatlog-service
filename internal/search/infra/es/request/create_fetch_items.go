package request

import (
	"github.com/aquasecurity/esquery"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/consts"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/query"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/request/agg"
	stdAgg "gitlab.trgdev.com/gotrg/white-label/services/catalog/pkg/es/agg_factory"
)

type m = map[string]any

/**
 * Builds request for fetching families.
 * 1. It creates query filters to form a basic set of products.
 * 2. It attaches aggregate functions to group a set of products into families.
 */
func CreateFetchItems(criteria *model.SearchCriteria) (m, error) {

	// Create request with filters to form basic set of products
	baseQuery, err := query.CreateQuery(criteria.GetFilters())
	if err != nil {
		return nil, err
	}

	// Create aggregation to group basic set of products into families
	aggs, err := createAggToFetchCollapsedFamilies(criteria)
	if err != nil {
		return nil, err
	}

	customQuery := m{
		"aggs": m{},
		// size should be 0 because we use aggregation
		// otherwise we will get items  by query
		"size":    0,
		"_source": false,
	}

	query, ok := baseQuery.Map()["bool"].(map[string]interface{})
	if ok && len(query) > 0 {
		// We needd only query property from parentRequest.
		// Conditions in query build base set of products.
		// Other properties are not needed because we use aggregation to fetch family properties
		// To sort families we use sorting in aggregation.
		customQuery["query"] = baseQuery.Map()
	}

	for _, agg := range aggs {
		customQuery["aggs"].(m)[agg.Name()] = agg.Map()
	}

	return customQuery, nil
}

func createAggToFetchCollapsedFamilies(
	criteria *model.SearchCriteria,
) ([]*esquery.CustomAggMap, error) {

	size := int(criteria.GetPageSize())
	from := int(criteria.GetPage()-1) * int(size)

	// this root aggregation to fetch families,
	// other aggregations will be added as child into this
	rootAgg := stdAgg.TermsAgg(
		"family.name",
		stdAgg.WithSize(consts.MAX_ALLOWED_FAMILY_QUANTITY),
		stdAgg.WithSubAggs(m{
			consts.BUCKET_NAME_FOR_OFFSET: m{
				// this is pipeline aggregation.
				// It is executed after the underlying aggregation dataset has been created.
				// So, size of the root aggregation should be greater than sum of the oconstsset and size.
				// And for the same reason we don't use the sorting here. Sorting placed  into the base aggregation.
				// @see https://www.elastic.co/guide/en/elasticsearch/reference/8.8/search-aggregations-pipeline-bucket-sort-aggregation.html
				"bucket_sort": m{
					"from": from,
					"size": size,
				},
			},
		}),
	)

	if criteria.GetOrderBy().GetField() == consts.FAMILY_KEY_FIELD {
		rootAgg["terms"].(m)["order"] = []any{
			m{"_key": criteria.GetOrderBy().GetDirection()},
		}
	}

	err := agg.AppendSubAggAndSorting(criteria, rootAgg)
	if err != nil {
		return nil, err
	}

	return []*esquery.CustomAggMap{
		// aggregation to fetch total count families
		// @see https://www.elastic.co/guide/en/elasticsearch/reference/8.8/search-aggregations-metrics-cardinality-aggregation.html
		esquery.CustomAgg("count", stdAgg.ValueAgg("family.name", "cardinality")),
		esquery.CustomAgg(consts.FAMILY_KEY_FIELD, rootAgg),
	}, nil
}
