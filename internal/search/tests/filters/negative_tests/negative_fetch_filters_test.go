package negative_tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "gitlab.trgdev.com/gotrg/white-label/services/catalog/api/catalog/proto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/dto"
	locationDto "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/tests"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/tests/mocks"
)

func TestSearchServer_Exec(t *testing.T) {

	t.Run("input request with invalid json", func(t *testing.T) {
		var err error

		os.Setenv("ES_HOST", "http://localhost:9200")
		os.Setenv("ES_INDEX", "prefix_")
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")
		defer dc.Reset()

		request := &pb.FiltersRequest{
			Body: "{ i am invalid json }",
		}

		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Filters(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Equal(t, "rpc error: code = InvalidArgument desc = invalid request body", err.Error())
	})

	t.Run("input request with unprocessable filter", func(t *testing.T) {
		var err error
		esHost := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Elastic-Product", "Elasticsearch")
			w.Header().Set("Content-Type", "application/json")
		}))
		defer esHost.Close()

		os.Setenv("ES_HOST", esHost.URL)

		os.Setenv("ES_INDEX", "prefix_")
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")
		defer dc.Reset()
		request := &pb.FiltersRequest{
			Body: `{
				"price": {
					"from": "100",
					"to": "invalid"
				}
			}`,
		}

		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Filters(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Equal(t, "rpc error: code = InvalidArgument desc = invalid argument: value, filter key = 'price', operation type = 'range'", err.Error())
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
		defer dc.Reset()

		request := &pb.FiltersRequest{}

		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Filters(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)
		assert.Equal(t, "rpc error: code = Internal desc = invalid response format from es", err.Error())
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
		defer dc.Reset()

		request := &pb.FiltersRequest{}

		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Filters(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Equal(t, "rpc error: code = Internal desc = invalid response format from es", err.Error())

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
		defer dc.Reset()

		request := &pb.FiltersRequest{}

		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Filters(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Equal(t, "rpc error: code = Internal desc = not found aggregations in es response", err.Error())

		dc.Reset()
	})

	t.Run("category options fetcher return error", func(t *testing.T) {

		esRespFixture := tests.ReadJsonFile("./fixtures/success_es_resp.json")

		host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// we set the header so that the client thinks
			// that the real response is from the server
			w.Header().Set("X-Elastic-Product", "Elasticsearch")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, esRespFixture)
		}))

		defer host.Close()

		os.Setenv("ES_HOST", host.URL)
		os.Setenv("ES_INDEX", "prefix_")
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")
		defer dc.Reset()

		request := &pb.FiltersRequest{}

		mockCategoryRepo := mocks.NewCategoryRepoMock()
		mockCategoryRepo.On("FindByIds", mock.Anything, mock.Anything).Return([]dto.Category{}, errors.New("findByIds return error"))
		mockLocationRepo := mocks.NewLocationRepoMock()
		mockLocationRepo.On("FindByIds", mock.Anything, mock.Anything).Return([]locationDto.Location{}, nil)

		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Filters(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Equal(t, "rpc error: code = Internal desc = internal error: findByIds return error", err.Error())

		dc.Reset()
	})

	t.Run("location options fetcher return error", func(t *testing.T) {

		esRespFixture := tests.ReadJsonFile("./fixtures/success_es_resp.json")

		host := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// we set the header so that the client thinks
			// that the real response is from the server
			w.Header().Set("X-Elastic-Product", "Elasticsearch")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, esRespFixture)
		}))

		defer host.Close()

		os.Setenv("ES_HOST", host.URL)
		os.Setenv("ES_INDEX", "prefix_")
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")
		defer dc.Reset()

		request := &pb.FiltersRequest{}

		mockCategoryRepo := mocks.NewCategoryRepoMock()
		mockCategoryRepo.On("FindByIds", mock.Anything, mock.Anything).Return([]dto.Category{}, nil)
		mockLocationRepo := mocks.NewLocationRepoMock()
		mockLocationRepo.On("FindByIds", mock.Anything, mock.Anything).Return([]locationDto.Location{}, errors.New("findByIds return error"))

		service := dc.SearchGrpcController.Use()
		actualResult, err := service.Filters(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, actualResult)

		assert.Equal(t, "rpc error: code = Internal desc = internal error: findByIds return error", err.Error())

		dc.Reset()
	})
}
