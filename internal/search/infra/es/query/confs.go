package query

import (
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/query/builders"
	i "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/query/iface"
)

var customBuilders = map[string]i.QueryBuilder{
	"location": builders.Location,
	"inStock":  builders.InStock,
}

var buildersByOperationType = map[string]i.QueryBuilder{
	model.TermsOp: builders.Terms,
	model.RangeOp: builders.Range,
	model.MatchOp: builders.Match,
}
