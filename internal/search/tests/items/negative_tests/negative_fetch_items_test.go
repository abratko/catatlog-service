package negative_tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "gitlab.trgdev.com/gotrg/white-label/services/catalog/api/catalog/proto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/tests/mocks"
)

func TestSearchServer_FetchItems(t *testing.T) {

	t.Run("input request with invalid json", func(t *testing.T) {
		var err error

		os.Setenv("ES_HOST", "http://localhost:9200")
		os.Setenv("ES_INDEX", "prefix_")
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")

		request := &pb.ItemsRequest{
			Body: "{ i am invalid json }",
		}

		mocks.NewCategoryRepoMock()
		mocks.NewLocationRepoMock()
		service := dc.SearchGrpcController.Use()

		actualResult, err := service.Items(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Equal(t, "rpc error: code = InvalidArgument desc = invalid request body", err.Error())

		dc.Reset()
	})

	t.Run("es send redirect", func(t *testing.T) {
		var err error
		esHost := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://localhost:9200/prefix_catalog_products/_search", http.StatusSeeOther)
		}))
		defer esHost.Close()

		os.Setenv("ES_HOST", esHost.URL)
		os.Setenv("ES_INDEX", "prefix_")
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")

		request := &pb.ItemsRequest{}

		mocks.NewCategoryRepoMock()
		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Items(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)
		assert.Equal(t, "rpc error: code = Internal desc = app service error: invalid response format from es", err.Error())

		dc.Reset()
	})

	t.Run("response from es does not contain valid json", func(t *testing.T) {
		host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// we set the header so that the client thinks
			// that the real response is from the server
			w.Header().Set("X-Elastic-Product", "Elasticsearch")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, "i am not valid json")
		}))

		defer host.Close()

		os.Setenv("ES_HOST", host.URL)
		os.Setenv("ES_INDEX", "prefix_")
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")

		request := &pb.ItemsRequest{}

		mocks.NewCategoryRepoMock()
		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Items(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Contains(t, err.Error(), "invalid response format from es")

		dc.Reset()
	})

	t.Run("response from es is JSON but contains invalid fields", func(t *testing.T) {
		host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// we set the header so that the client thinks
			// that the real response is from the server
			w.Header().Set("X-Elastic-Product", "Elasticsearch")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"unexpected_field": "unexpected_value"}`)
		}))

		defer host.Close()

		os.Setenv("ES_HOST", host.URL)
		os.Setenv("ES_INDEX", "prefix_")
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")

		request := &pb.ItemsRequest{}

		mocks.NewCategoryRepoMock()
		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Items(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Contains(t, err.Error(), "invalid format for totalCount")

		dc.Reset()
	})
}
