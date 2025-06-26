package dc

import (
	"context"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"gitlab.com/gorib/env"

	"github.com/abratko/dc"
	d "github.com/abratko/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/pkg/es"
)

var EsClient = dc.Provider(
	func() *elasticsearch.Client {

		esHost := env.NeedValue[string]("ES_HOST")
		config := elasticsearch.Config{
			Addresses:            []string{esHost},
			UseResponseCheckOnly: true,
		}

		if env.Value("ES_DEBUG", false) {
			config.Logger = &estransport.ColorLogger{
				EnableRequestBody:  true,
				EnableResponseBody: true,
				Output:             os.Stdout,
			}
		}

		es, err := elasticsearch.NewClient(config)
		if err != nil {
			panic(err)
		}

		return es
	},
)

var ProductReader = dc.Provider(
	func() *es.ReadAdapter {
		return es.NewReadAdapter(
			EsClient.Use(),
			env.NeedValue[string]("ES_INDEX")+"product",
		)
	},
)
var ReadProductDb = d.Provider(
	func() func(context.Context, map[string]any) (*esapi.Response, error) {
		return ProductReader.Use().Read
	})
