package request

import (
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/query"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/request/agg"
)

func CreateFetchFilterOpts(
	filters *model.Filters,
) (map[string]any, error) {
	baseQuery, err := query.CreateQuery(filters, filters.OnlyFacetedKeys()...)
	if err != nil {
		return nil, err
	}

	facetAggregations := make(map[string]interface{})
	// Values for all faceted filters on client side should be updated after request.
	// So, we need iterate over all faceted filters and create aggs for each filter
	// to get options including current filter value.
	for _, facetedFilter := range filters.OnlyFaceted() {
		// some agg can depend from values of other input filters.
		// For example, location agg can depend from inStock filter.
		// we don`t know from which filters, so we need to pass all filters as args
		agg, err := agg.CreateAggs(facetedFilter, filters)
		if err != nil {
			return nil, err
		}

		// For subquery we exclude not faceted filters because we already use them in base query.
		// We also exclude current faceted filter because
		// We want to get options including current filter value.
		excludedKeys := append(filters.OnlyNotFacetedKeys(), facetedFilter.Key())
		subquery, err := query.CreateQuery(filters, excludedKeys...)
		if err != nil {
			return nil, err
		}

		isAggEmpty := len(subquery.Map()["bool"].(map[string]any)) == 0
		if isAggEmpty {
			facetAggregations[facetedFilter.Key()] = agg

			continue
		}

		facetAggregations[facetedFilter.Key()] = map[string]any{
			"filter": subquery.Map(),
			"aggs": map[string]any{
				facetedFilter.Key(): agg,
			},
		}

	}

	request := map[string]any{
		"aggs":    facetAggregations,
		"size":    0,
		"_source": false,
	}

	query := baseQuery.Map()
	if len(query["bool"].(map[string]any)) > 0 {
		request["query"] = query
	}

	return request, nil
}
