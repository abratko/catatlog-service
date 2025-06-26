package positive

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/errs"
)

func TestLocationRepo_InitFailed(t *testing.T) {

	t.Run("FindByIds returns ErrFetchItems and not reinit", func(t *testing.T) {
		var countCalls int
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			countCalls++
			// we set the header so that the client thinks
			// that the real response is from the server
			w.Header().Set("X-Elastic-Product", "Elasticsearch")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, "i am not valid json")
		}))

		os.Setenv("ES_HOST", ts.URL)
		os.Setenv("ES_INDEX", "prefix_")
		os.Setenv("ES_DEBUG", "false")
		os.Setenv("CACHE_EXP_AFTER", fmt.Sprintf("%v", 1*time.Second))
		defer os.Unsetenv("ES_HOST")
		defer os.Unsetenv("ES_INDEX")
		defer dc.Reset()
		defer ts.Close()

		repo := dc.CategoryRepo.Use()

		result1, err := repo.FindByIds(context.Background(), []string{"1"})

		time.Sleep(2 * time.Second)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrFetchItems))
		assert.Nil(t, result1)
		assert.Equal(t, 1, countCalls)
	})
}
