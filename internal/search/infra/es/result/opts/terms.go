package opts

import (
	"fmt"
	"reflect"

	"github.com/Jeffail/gabs/v2"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

func Terms(filterKey string, buckets []*gabs.Container) ([]model.TermOpt, error) {
	if len(buckets) == 0 {
		return nil, fmt.Errorf("not found buckets in response for filter key: %s", filterKey)
	}

	terms := make([]model.TermOpt, 0, len(buckets))
	for _, bucket := range buckets {
		value := bucket.Path("key").Data()
		kind := reflect.TypeOf(value).Kind()

		if kind != reflect.String {
			value = fmt.Sprintf("%v", value)
		}

		term := model.TermOpt{
			Value: value.(string),
			Count: int(bucket.Path("doc_count").Data().(float64)),
		}
		terms = append(terms, term)
	}

	return terms, nil
}
