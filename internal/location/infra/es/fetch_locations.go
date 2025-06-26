package es

import (
	"bufio"
	"context"
	"errors"

	"github.com/Jeffail/gabs/v2"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/pkg/es"
)

type Fetcher struct {
	db *es.ReadAdapter
}

func NewFetcher(
	db *es.ReadAdapter,
) *Fetcher {
	return &Fetcher{
		db: db,
	}
}
func (c *Fetcher) Exec(
	ctx context.Context,
) (chan dto.Location, error) {

	res, err := c.db.Read(ctx, map[string]any{"size": 1000})
	if err != nil {
		return nil, err
	}

	return c.createLocationsGenerator(res)
}

func (c *Fetcher) createLocationsGenerator(res *esapi.Response) (chan dto.Location, error) {
	reader := bufio.NewReader(res.Body)

	parsedJSON, err := gabs.ParseJSONBuffer(reader)
	if err != nil {
		return nil, errors.New("invalid response format from es")
	}

	hits := parsedJSON.Path("hits.hits").Children()
	if len(hits) == 0 {
		return nil, errors.New("not found hits in es response")
	}

	ch := make(chan dto.Location)
	go func() {
		// Закрываем канал после завершения отправки данных
		defer close(ch)

		for _, hit := range hits {

			if !hit.ExistsP("_source") {
				continue
			}

			location := dto.Location{
				Id:      int(hit.Path("_source.id").Data().(float64)),
				ZipCode: hit.Path("_source.zip_code").Data().(string),
				Country: dto.Country{
					Name: hit.Path("_source.country.name").Data().(string),
					Code: hit.Path("_source.country.code").Data().(string),
				},
				State: dto.State{
					Name: hit.Path("_source.state.name").Data().(string),
					Code: hit.Path("_source.state.code").Data().(string),
				},
				City: hit.Path("_source.city").Data().(string),
				Address: dto.Address{
					Line1: hit.Path("_source.address.line1").Data().(string),
					Line2: hit.Path("_source.address.line2").Data().(string),
					Line3: hit.Path("_source.address.line3").Data().(string),
				},
				IsResidential: hit.Path("_source.is_residential").Data().(bool),
				Phone:         hit.Path("_source.phone").Data().(string),
				Name:          hit.Path("_source.name").Data().(string),
			}

			ch <- location
		}
	}()

	return ch, nil
}
