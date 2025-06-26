package result

import (
	"bufio"
	"errors"

	"github.com/Jeffail/gabs/v2"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/result/opts"
)

type valsByFilterKey map[string]any

func CreateOptDtos(esResponse *esapi.Response) (valsByFilterKey, error) {
	var err error

	reader := bufio.NewReader(esResponse.Body)

	parsedJSON, err := gabs.ParseJSONBuffer(reader)
	if err != nil {
		return nil, errors.New("invalid response format from es")
	}

	aggregations := parsedJSON.Path("aggregations").ChildrenMap()
	if len(aggregations) == 0 {
		return nil, errors.New("not found aggregations in es response")
	}

	vals := make(map[string]any)
	for filterKey, aggregation := range aggregations {
		values, err := opts.Create(filterKey, aggregation)
		if err != nil {
			return nil, err
		}

		vals[filterKey] = values
	}

	return vals, nil
}
