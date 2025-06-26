package config

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type EsReader interface {
	Read(ctx context.Context, queryDto map[string]any) (*esapi.Response, error)
}

type EsClientProvider = func() EsReader

var esClientProvider EsClientProvider

func Init(esClient EsClientProvider) {
	esClientProvider = esClient
}

func GetEsClient() EsReader {
	return esClientProvider()
}
