package es

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs/v2"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// type Search func(o ...func(*esapi.SearchRequest)) (*esapi.Response, error)
// type Client interface {
// 	Search(o ...func(*esapi.SearchRequest)) (*esapi.Response, error)
// 	// func (f Search) WithContext(v context.Context) func(*esapi.SearchRequest)
// }

// type func (f Client.Search) WithContext(v context.Context) func(*esapi.SearchRequest)

type ReadAdapter struct {
	client *elasticsearch.Client
	index  string
}

func NewReadAdapter(
	client *elasticsearch.Client,
	index string,
) *ReadAdapter {
	return &ReadAdapter{
		client: client,
		index:  index,
	}
}

// I don't know how write tests for this function
// because i don't know how to mock elasticsearch.Client
// It is not an interface, it is a function!!! with methods
func (r *ReadAdapter) Read(
	ctx context.Context,
	queryDto map[string]any,
) (*esapi.Response, error) {
	json, err := json.Marshal(queryDto)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write([]byte(json))
	qRes, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, err
	}

	if qRes.StatusCode == 200 {
		return qRes, nil
	}

	reader := bufio.NewReader(qRes.Body)

	parsedJSON, err := gabs.ParseJSONBuffer(reader)
	if err != nil {
		return nil, errors.New("invalid response format from es")
	}

	errorType, ok := parsedJSON.Path("error.type").Data().(string)
	if !ok {
		return nil, errors.New("unknown error by reading from es")
	}

	return nil, errors.New(errorType)
}
