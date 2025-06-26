package agg_factory

func ValueAgg(field string, operation string) m {
	return m{
		operation: m{
			"field": field,
		},
	}
}
