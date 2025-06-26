package opts_test

import (
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/stretchr/testify/assert"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/result/opts"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name             string
		filterKey        string
		fieldAggregation *gabs.Container
		expectedResult   any
		expectedError    bool
		errorContains    string
	}{
		{
			name:      "terms aggregation with buckets at root",
			filterKey: "category",
			fieldAggregation: gabs.Wrap(map[string]any{
				"buckets": []any{
					map[string]any{
						"key":       "electronics",
						"doc_count": 10.0,
					},
					map[string]any{
						"key":       "books",
						"doc_count": 5.0,
					},
				},
			}),
			expectedResult: []model.TermOpt{
				{Value: "electronics", Count: 10},
				{Value: "books", Count: 5},
			},
			expectedError: false,
		},
		{
			name:      "terms aggregation with nested buckets",
			filterKey: "category",
			fieldAggregation: gabs.Wrap(map[string]any{
				"category": map[string]any{
					"buckets": []any{
						map[string]any{
							"key":       "electronics",
							"doc_count": 10.0,
						},
						map[string]any{
							"key":       "books",
							"doc_count": 5.0,
						},
					},
				},
			}),
			expectedResult: []model.TermOpt{
				{Value: "electronics", Count: 10},
				{Value: "books", Count: 5},
			},
			expectedError: false,
		},
		{
			name:      "range aggregation at root",
			filterKey: "price",
			fieldAggregation: gabs.Wrap(map[string]any{
				"min": 10.0,
				"max": 100.0,
			}),
			expectedResult: model.RangeOpt{
				From: 10.0,
				To:   100.0,
			},
			expectedError: false,
		},
		{
			name:      "range aggregation with nested path",
			filterKey: "price",
			fieldAggregation: gabs.Wrap(map[string]any{
				"price": map[string]any{
					"min": 10.0,
					"max": 100.0,
				},
			}),
			expectedResult: model.RangeOpt{
				From: 10.0,
				To:   100.0,
			},
			expectedError: false,
		},
		{
			name:      "custom factory for location with nested buckets",
			filterKey: "location",
			fieldAggregation: gabs.Wrap(map[string]any{
				"location": map[string]any{
					"location": map[string]any{
						"location": map[string]any{
							"buckets": []any{
								map[string]any{
									"key":       "US",
									"doc_count": 10.0,
								},
							},
						},
					},
				},
			}),
			expectedResult: []model.TermOpt{
				{Value: "US", Count: 10},
			},
			expectedError: false,
		},
		{
			name:      "custom factory for location as root",
			filterKey: "location",
			fieldAggregation: gabs.Wrap(map[string]any{
				"location": map[string]any{
					"location": map[string]any{
						"buckets": []any{
							map[string]any{
								"key":       "US",
								"doc_count": 10.0,
							},
						},
					},
				},
			}),
			expectedResult: []model.TermOpt{
				{Value: "US", Count: 10},
			},
			expectedError: false,
		},
		{
			name:      "no matching aggregation type",
			filterKey: "unknown",
			fieldAggregation: gabs.Wrap(map[string]any{
				"unknown": map[string]any{
					"some_field": "value",
				},
			}),
			expectedError: true,
			errorContains: "not found values factory by filter key: unknown",
		},
		{
			name:      "invalid range aggregation - missing min",
			filterKey: "price",
			fieldAggregation: gabs.Wrap(map[string]any{
				"max": 100.0,
			}),
			expectedError: true,
			errorContains: "not found values factory by filter key: price",
		},
		{
			name:      "invalid range aggregation - missing max",
			filterKey: "price",
			fieldAggregation: gabs.Wrap(map[string]any{
				"min": 10.0,
			}),
			expectedError: true,
			errorContains: "not found max aggregation response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := opts.Create(tt.filterKey, tt.fieldAggregation)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
