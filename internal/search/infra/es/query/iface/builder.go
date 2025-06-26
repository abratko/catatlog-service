package iface

import (
	"github.com/aquasecurity/esquery"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

type QueryBuilder func(
	filter *model.Filter,
	filters *model.Filters,
	query *esquery.BoolQuery,
) bool
