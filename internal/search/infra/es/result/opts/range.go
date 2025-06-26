package opts

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

func Range(filterKey string, aggregation *gabs.Container) (model.RangeOpt, error) {
	if !aggregation.ExistsP("min") {
		return model.RangeOpt{}, fmt.Errorf("not found min aggregation response for filter key: %s", filterKey)
	}

	if !aggregation.ExistsP("max") {
		return model.RangeOpt{}, fmt.Errorf("not found max aggregation response for filter key: %s", filterKey)
	}

	from := aggregation.Path("min").Data().(float64)
	to := aggregation.Path("max").Data().(float64)

	return model.RangeOpt{
		From: from,
		To:   to,
	}, nil
}
