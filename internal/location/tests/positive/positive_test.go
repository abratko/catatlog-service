package positive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/tests"
)

func TestLocationRepo_FindByIds(t *testing.T) {
	// Create test locations
	locations := []dto.Location{
		{
			Id:      1,
			ZipCode: "12345",
			Country: dto.Country{
				Name: "United States",
				Code: "US",
			},
			State: dto.State{
				Name: "California",
				Code: "CA",
			},
			City: "San Francisco",
			Address: dto.Address{
				Line1: "123 Main St",
				Line2: "Suite 100",
				Line3: "",
			},
			IsResidential: false,
			Name:          "Office",
			Phone:         "555-0123",
		},
		{
			Id:      2,
			ZipCode: "67890",
			Country: dto.Country{
				Name: "United States",
				Code: "US",
			},
			State: dto.State{
				Name: "New York",
				Code: "NY",
			},
			City: "New York",
			Address: dto.Address{
				Line1: "456 Park Ave",
				Line2: "",
				Line3: "",
			},
			IsResidential: true,
			Name:          "Home",
			Phone:         "555-0124",
		},
	}

	// Create a mock fetch function that returns our test locations
	mockFetch := func(ctx context.Context) (chan dto.Location, error) {
		ch := make(chan dto.Location, len(locations))
		go func() {
			defer close(ch)
			for _, loc := range locations {
				ch <- loc
			}
		}()
		return ch, nil
	}

	// Create the repository with a shorter timeout for testing
	repo := app.NewRepo(mockFetch, "3m")

	// Wait for initialization to complete
	time.Sleep(100 * time.Millisecond)

	tests := []struct {
		name    string
		ids     []string
		want    int
		wantErr bool
	}{
		{
			name: "find single location",
			ids:  []string{"1"},
			want: 1,
		},
		{
			name: "find multiple locations",
			ids:  []string{"1", "2"},
			want: 2,
		},
		{
			name: "find non-existent location",
			ids:  []string{"999"},
			want: 0,
		},
		{
			name: "find empty list",
			ids:  []string{},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByIds(context.Background(), tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("FindByIds() got %v locations, want %v", len(got), tt.want)
			}
		})
	}
}

func TestLocationRepo_Create(t *testing.T) {

	t.Run("init during creation and reinit after expiration", func(t *testing.T) {
		var actualEsRequestBody []byte

		fixtureEsResponse := tests.ReadJsonFile("fixtures/fetch_all_es_resp_init_call.json")
		fixtureEsResponseReinit := tests.ReadJsonFile("fixtures/fetch_all_es_resp_reinit_call.json")
		expectedEsRequest := tests.ReadJsonFile("expected/fetch_all_es_req.json")
		expectedResult1 := tests.ReadJsonFile("expected/result1.json")
		expectedResult2 := tests.ReadJsonFile("expected/result2.json")
		expEsReqPath := "/prefix_facility/_search"

		var actPath string
		var countCalls int
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

			if countCalls == 0 {
				fmt.Fprintln(w, fixtureEsResponse)
			} else {
				fmt.Fprintln(w, fixtureEsResponseReinit)
			}

			countCalls++
		}))

		os.Setenv("ES_HOST", ts.URL)
		os.Setenv("ES_INDEX", "prefix_")
		os.Setenv("ES_DEBUG", "false")
		os.Setenv("CACHE_EXP_AFTER", fmt.Sprintf("%v", 3*time.Second))
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")
		defer dc.Reset()
		defer ts.Close()

		repo := dc.LocationRepo.Use()

		result1, err := repo.FindByIds(context.Background(), []string{"1"})
		assert.NoError(t, err)
		assert.NotNil(t, result1)
		result1AsJson, _ := json.Marshal(result1)

		time.Sleep(4 * time.Second)
		result2, err := repo.FindByIds(context.Background(), []string{"12"})
		result2AsJson, _ := json.Marshal(result2)
		repo.CancelReinit()

		assert.NoError(t, err)
		assert.NotNil(t, result2)

		assert.Equal(t, 2, countCalls)
		assert.Equal(t, expEsReqPath, actPath)

		assert.Equal(t, expectedEsRequest, string(actualEsRequestBody))

		assert.Equal(t, expectedResult1, string(result1AsJson))
		assert.Equal(t, expectedResult2, string(result2AsJson))
	})
}
