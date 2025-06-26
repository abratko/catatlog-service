package dc

import (
	d "github.com/abratko/dc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/service"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/service/opt_fetcher"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/contr/grpc"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es"
)

var SearchGrpcController = d.Provider(func() *grpc.Controller {
	var optsFetchers = []service.OptFetcher{
		opt_fetcher.NewCategoryOptFetcher(CategoryRepo.Use().FindByIds),
		opt_fetcher.NewLocationOptFetcher(LocationRepo.Use().FindByIds),
	}

	var filtersFactory = service.NewFiltersFactory(optsFetchers)
	var esFetchFilters = es.NewFetchFilters(filtersFactory.Exec)
	var esFetchFilterOpts = es.NewFetchFilterOpts(ReadProductDb.Use())

	var appFetchFiltersWithOpts = app.NewFetchFiltersWithOpts(
		esFetchFilters.Exec,
		esFetchFilterOpts.Exec,
	)

	var esFetchItems = es.NewFetchItems(ReadProductDb.Use())
	var appFetchItems = app.NewFetchItems(
		esFetchFilters.Exec,
		esFetchItems.Exec,
	)

	return grpc.NewController(
		appFetchFiltersWithOpts.Exec,
		appFetchItems.Exec,
	)
})
