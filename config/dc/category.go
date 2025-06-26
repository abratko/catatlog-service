package dc

import (
	"context"
	"fmt"
	"time"

	"github.com/abratko/dc"
	"gitlab.com/gorib/env"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/app/dto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/category/infra/es"
	pkgEs "gitlab.trgdev.com/gotrg/white-label/services/catalog/pkg/es"
)

type CatRepo interface {
	FindByIds(ctx context.Context, ids []string) ([]dto.Category, error)
	CancelReinit()
}

var CategoryRepo = dc.Provider(
	func() CatRepo {
		var index = env.NeedValue[string]("ES_INDEX") + "category"
		var duration = env.Value("CACHE_EXP_AFTER", fmt.Sprintf("%v", 10*time.Minute))

		var esAdapter = pkgEs.NewReadAdapter(EsClient.Use(), index)

		var fetchCategories = es.NewFetchCategories(esAdapter)

		return app.NewCategoryRepo(fetchCategories.Exec, duration)
	})
