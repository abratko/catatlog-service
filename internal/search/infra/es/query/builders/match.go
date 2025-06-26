package builders

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/aquasecurity/esquery"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
)

func Match(
	filter *model.Filter,
	filters *model.Filters,
	query *esquery.BoolQuery,
) bool {
	values := filter.Value()
	kindValues := reflect.TypeOf(values).Kind()
	if kindValues != reflect.String {
		return false
	}

	// Replace all not numeric, alphabetical and space symbols
	re := regexp.MustCompile(`[^\dA-Za-z\s]+`)
	cleanedValue := re.ReplaceAllString(values.(string), " ")

	if cleanedValue == "" {
		return false
	}

	// searchParams := map[string]any{
	// 	"query":                fmt.Sprintf("(%1$s) OR (%1$s*)", cleanedValue),
	// 	"fields":               strings.Split(fieldName, ","),
	// 	"default_operator":     "AND",
	// 	"fuzziness":            2,
	// 	"fuzzy_max_expansions": 3,
	// }
	// this is alternative structure previosly used in the code
	// It was used to create query_string search
	// @see https://www.elastic.co/guide/en/elasticsearch/reference/6.8/query-dsl-query-string-query.html#_multi_field
	//
	// esquery does not support query_string search, so we need to use MultiMatch instead

	// filedName = "name^5,description,sku"
	query.Must(
		esquery.MultiMatch().
			Query(cleanedValue).
			Fields(strings.Split(filter.Field(), ",")...).
			Operator(esquery.OperatorAnd).
			Fuzziness("2").
			MaxExpansions(3),
	)

	return true
}
