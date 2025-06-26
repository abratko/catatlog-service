package agg_factory

type m = map[string]any

type TermsAggOption func(termsAgg m)

func WithSubAggs(subAggs m) TermsAggOption {
	return func(termsAgg m) {
		termsAgg["aggs"] = subAggs
	}
}

func WithOrder(order []interface{}) TermsAggOption {
	return func(termsAgg m) {
		termsAgg["order"] = order
	}
}

func WithSize(size int) TermsAggOption {
	return func(termsAgg m) {
		termsAgg["terms"].(m)["size"] = size
	}
}

// termsAgg создает новый TermsAgg с заданными параметрами
func TermsAgg(
	field string,
	opts ...TermsAggOption,
) m {
	agg := m{
		"terms": m{
			"field": field,
		},
	}

	for _, opt := range opts {
		opt(agg)
	}

	return agg
}
