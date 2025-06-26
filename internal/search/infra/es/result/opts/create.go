package opts

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

var factoriesByKey = map[string]func(string, *gabs.Container) (any, error){
	"location": Location,
}

func Create(filterKey string, fieldAggregation *gabs.Container) (any, error) {
	if fieldAggregation.ExistsP("buckets") {
		return Terms(filterKey, fieldAggregation.Path("buckets").Children())
	}

	if fieldAggregation.ExistsP(filterKey + ".buckets") {
		return Terms(filterKey, fieldAggregation.Path(filterKey+".buckets").Children())
	}

	if fieldAggregation.ExistsP("min") {
		return Range(filterKey, fieldAggregation)
	}

	if fieldAggregation.ExistsP(filterKey + ".min") {
		return Range(filterKey, fieldAggregation.Path(filterKey))
	}

	if factory, found := factoriesByKey[filterKey]; found {
		return factory(filterKey, fieldAggregation)
	}

	return nil, fmt.Errorf("not found values factory by filter key: %s", filterKey)
}
