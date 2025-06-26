package dc

import (
	"context"
	"fmt"
	"time"

	"github.com/abratko/dc"
	"gitlab.com/gorib/env"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/location/infra/es"
	pkgEs "gitlab.trgdev.com/gotrg/white-label/services/catalog/pkg/es"
)

type LocRepo interface {
	FindByIds(ctx context.Context, ids []string) ([]dto.Location, error)
	CancelReinit()
}

var LocationRepo = dc.Provider(
	func() LocRepo {
		var index = env.NeedValue[string]("ES_INDEX") + "facility"
		var duration = env.Value("CACHE_EXP_AFTER", fmt.Sprintf("%v", 10*time.Minute))

		var esAdapter = pkgEs.NewReadAdapter(EsClient.Use(), index)
		var fetcher = es.NewFetcher(esAdapter)

		return app.NewRepo(
			fetcher.Exec,
			duration,
		)
	})
