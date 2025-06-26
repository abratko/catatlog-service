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

	pb "gitlab.trgdev.com/gotrg/white-label/services/catalog/api/catalog/proto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/tests"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/tests/mocks"
)

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
	expGrpcResp := &pb.ItemsResult{}

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

	defer ts.Close()

	os.Setenv("ES_HOST", ts.URL)
	os.Setenv("ES_INDEX", "prefix_")
	os.Setenv("ES_DEBUG", "false")
	defer os.Unsetenv("ES_HOST")
	defer os.Unsetenv("ES_INDEX")
	defer dc.Reset()

	request := &pb.ItemsRequest{
		Body: testDto.fixtureRequest,
	}

	mocks.NewCategoryRepoMock()
	mocks.NewLocationRepoMock()
	service := dc.SearchGrpcController.Use()
	actualResult, err := service.Items(context.Background(), request)

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

func TestSearchServer_FetchItems(t *testing.T) {

	esResponseWithoutSortingField := tests.ReadJsonFile("./fixtures/es_resp_without_sorting_field.json")
	sharedExpGrpcRespAsJson := tests.ReadJsonFile("./expected/grpc_resp.json")
	expEsRequestPath := "/prefix_product/_search"

	var positiveTests = []PositiveTestDto{
		{
			name:              "request is empty",
			fixtureRequest:    "",
			fixtureEsResponse: esResponseWithoutSortingField,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_when_all_filters_empty.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "all filters set and order by name",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_all_filters_set_order_by_name.json"),
			fixtureEsResponse: esResponseWithoutSortingField,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_all_filters_set_order_by_name.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "all filters set and order by lowestAsk",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_all_filters_set_order_by_lowest_ask.json"),
			fixtureEsResponse: esResponseWithoutSortingField,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_all_filters_set_order_by_lowest_ask.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "all filters set and order by expectedProfit",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_all_filters_set_order_by_expected_profit.json"),
			fixtureEsResponse: tests.ReadJsonFile("./fixtures/es_resp_with_expected_profit.json"),

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_all_filters_set_order_by_expected_profit.json"),
			expGrpcRespAsJson: tests.ReadJsonFile("./expected/grpc_resp_with_expected_profit.json"),
		},
		{
			name:              "with part filters only",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_part_filters_only.json"),
			fixtureEsResponse: esResponseWithoutSortingField,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_part_filters_only.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "with order by only",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_with_order_by_only.json"),
			fixtureEsResponse: esResponseWithoutSortingField,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_with_order_by_only.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "with pagination only",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_with_pagination_only.json"),
			fixtureEsResponse: esResponseWithoutSortingField,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_with_pagination_only.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
		{
			name:              "skip all unknown filters",
			fixtureRequest:    tests.ReadJsonFile("./fixtures/req_with_unknown_filters_set.json"),
			fixtureEsResponse: esResponseWithoutSortingField,

			expEsRequestPath:  expEsRequestPath,
			expEsRequest:      tests.ReadJsonFile("./expected/es_req_part_filters_only.json"),
			expGrpcRespAsJson: sharedExpGrpcRespAsJson,
		},
	}

	for _, testDto := range positiveTests {
		t.Run(testDto.name, func(t *testing.T) {
			runTest(t, testDto)
		})
	}
}
