package builders_test

// import (
// 	"testing"

// 	"github.com/aquasecurity/esquery"
// 	"github.com/stretchr/testify/assert"
// 	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/query_factory/builders"
// )

// func TestTerms(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fieldName string
// 		values    any
// 		expected  *esquery.BoolQuery
// 	}{
// 		{
// 			name:      "valid string slice",
// 			fieldName: "category",
// 			values:    []string{"books", "electronics"},
// 			expected: esquery.Bool().Filter(
// 				esquery.Terms("category", []string{"books", "electronics"}),
// 			),
// 		},
// 		{
// 			name:      "empty slice",
// 			fieldName: "category",
// 			values:    []string{},
// 			expected:  esquery.Bool(),
// 		},
// 		{
// 			name:      "invalid type",
// 			fieldName: "category",
// 			values:    "not a slice",
// 			expected:  esquery.Bool(),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			query := esquery.Bool()
// 			builders.Terms(tt.fieldName, tt.values, map[string]any{}, query)
// 			assert.Equal(t, tt.expected, query)
// 		})
// 	}
// }

// func TestMatch(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fieldName string
// 		values    any
// 		expected  *esquery.BoolQuery
// 	}{
// 		{
// 			name:      "valid string",
// 			fieldName: "name,description",
// 			values:    "test product",
// 			expected: esquery.Bool().Must(
// 				esquery.MultiMatch().
// 					Query("test product").
// 					Fields("name", "description").
// 					Operator(esquery.OperatorAnd).
// 					Fuzziness("2").
// 					MaxExpansions(3),
// 			),
// 		},
// 		{
// 			name:      "special characters",
// 			fieldName: "description",
// 			values:    "test@#$%^&*product",
// 			expected: esquery.Bool().Must(
// 				esquery.MultiMatch().
// 					Query("test product").
// 					Fields("description").
// 					Operator(esquery.OperatorAnd).
// 					Fuzziness("2").
// 					MaxExpansions(3),
// 			),
// 		},
// 		{
// 			name:      "invalid type",
// 			fieldName: "description",
// 			values:    123,
// 			expected:  esquery.Bool(),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			query := esquery.Bool()
// 			builders.Match(tt.fieldName, tt.values, map[string]any{}, query)
// 			assert.Equal(t, tt.expected, query)
// 		})
// 	}
// }

// func TestRange(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		fieldName string
// 		values    any
// 		expected  *esquery.BoolQuery
// 	}{
// 		{
// 			name:      "valid range",
// 			fieldName: "price",
// 			values: map[string]any{
// 				"from": 10,
// 				"to":   20,
// 			},
// 			expected: esquery.Bool().Filter(
// 				esquery.Range("price").Gte(10).Lte(20),
// 			),
// 		},
// 		{
// 			name:      "only from",
// 			fieldName: "price",
// 			values: map[string]any{
// 				"from": 10,
// 			},
// 			expected: esquery.Bool().Filter(
// 				esquery.Range("price").Gte(10),
// 			),
// 		},
// 		{
// 			name:      "only to",
// 			fieldName: "price",
// 			values: map[string]any{
// 				"to": 20,
// 			},
// 			expected: esquery.Bool().Filter(
// 				esquery.Range("price").Lte(20),
// 			),
// 		},
// 		{
// 			name:      "equal values",
// 			fieldName: "price",
// 			values: map[string]any{
// 				"from": 10,
// 				"to":   10,
// 			},
// 			expected: esquery.Bool(),
// 		},
// 		{
// 			name:      "value as strings are inserted as is",
// 			fieldName: "price",
// 			values: map[string]any{
// 				"from": "some value 1",
// 				"to":   "some value 2",
// 			},
// 			expected: esquery.Bool().Filter(
// 				esquery.Range("price").Gte("some value 1").Lte("some value 2"),
// 			),
// 		},
// 		{
// 			name:      "invalid type",
// 			fieldName: "price",
// 			values:    "not a map",
// 			expected:  esquery.Bool(),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			query := esquery.Bool()
// 			builders.Range(tt.fieldName, tt.values, map[string]any{}, query)
// 			assert.Equal(t, tt.expected, query)
// 		})
// 	}
// }
