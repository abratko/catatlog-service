package opts

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

func Location(filterKey string, aggregation *gabs.Container) (any, error) {
	path := "location.location.buckets"
	if aggregation.ExistsP(path) {
		return Terms(filterKey, aggregation.Path(path).Children())
	}

	if aggregation.ExistsP(filterKey + ".location.location.buckets") {
		return Terms(filterKey, aggregation.Path(filterKey+".location.location.buckets").Children())
	}

	return nil, fmt.Errorf("not found location aggregation in es response")
}
