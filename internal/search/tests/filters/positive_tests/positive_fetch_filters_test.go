package positive_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "gitlab.trgdev.com/gotrg/white-label/services/catalog/api/catalog/proto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/tests"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/tests/filters/positive_tests/fixtures"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/tests/mocks"
)

type MockCategoryRepo struct {
	mock.Mock
}

type PositiveTestDto struct {
	name              string
	fixtureRequest    string
	fixtureEsResponse string
	expEsRequestPath  string
	expEsRequest      string
	expGrpcRespAsJson string
}

func runTest(t *testing.T, testDto PositiveTestDto) {
	var actualEsRequestBody []byte
	expGrpcResp := &pb.FiltersResult{}

	err := json.Unmarshal([]byte(string(testDto.expGrpcRespAsJson)), expGrpcResp)
	if err != nil {
		t.Fatal("Ошибка при десериализации expGrpcRespAsJson:", err)
	}

	var actPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actPath = r.URL.Path

		var err error

		actualEsRequestBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		// we set the header so that the client thinks
		// that the real response is from the server
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, testDto.fixtureEsResponse)
	}))

	defer func() {
		ts.Close()
	}()

	os.Setenv("ES_HOST", ts.URL)
	os.Setenv("ES_INDEX", "prefix_")
	os.Setenv("ES_DEBUG", "false")
	defer os.Unsetenv("ES_HOST")
	defer os.Unsetenv("ES_INDEX")

	request := &pb.FiltersRequest{
		Body: testDto.fixtureRequest,
	}

	mocks.
		NewCategoryRepoMock().
		On("FindByIds", mock.Anything, mock.Anything).Return(fixtures.CategoriesFixture, nil)
	mocks.
		NewLocationRepoMock().
		On("FindByIds", mock.Anything, mock.Anything).Return(fixtures.LocationsFixture, nil)

	defer dc.Reset()

	service := dc.SearchGrpcController.Use()

	actualResult, err := service.Filters(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, actualResult)

	assert.Equal(t, string(testDto.expEsRequest), string(actualEsRequestBody))
	assert.Equal(t, testDto.expEsRequestPath, actPath)

	actualResultAsJson, err := json.Marshal(actualResult)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(testDto.expGrpcRespAsJson), string(actualResultAsJson))
}

func TestSearchServer_FetchFilters(t *testing.T) {

	sharedSuccessEsResponse := tests.ReadJsonFile("./fixtures/success_es_resp.json")
	sharedExpGrpcRespAsJson := tests.ReadJsonFile("./expected/grpc_resp.json")
	expEsRequestPath := "/prefix_product/_search"

	var positiveTests = []PositiveTestDto{
		{
			name:              "request is empty",
			fixtureRequest:    "",
			fixtureEsResponse: sharedSuccessEsResponse,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_when_all_filters_empty.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "all filters set",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_all_filters_set.json"),
			fixtureEsResponse: sharedSuccessEsResponse,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_when_all_filters_set.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "only not faceted filters set",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_only_not_faceted_filters_set.json"),
			fixtureEsResponse: sharedSuccessEsResponse,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_only_not_faceted_filters_set.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "only faceted filters set",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_only_faceted_filters_set.json"),
			fixtureEsResponse: sharedSuccessEsResponse,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_only_faceted_filters_set.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "only part of filters set",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_part_of_filters_set.json"),
			fixtureEsResponse: sharedSuccessEsResponse,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_part_of_filters_set.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "skip all unknown filters",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_with_unknown_filters_set.json"),
			fixtureEsResponse: sharedSuccessEsResponse,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_part_of_filters_set.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
	}

	for _, testDto := range positiveTests {
		t.Run(testDto.name, func(t *testing.T) {
			runTest(t, testDto)
		})
	}
}
