package model_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/errs"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
)

func TestTermsFilter(t *testing.T) {
	errorMsgForSetValue := fmt.Errorf("%w: value, filter key = 'category', operation type = 'terms'", errs.ErrInvalidArg)

	tests := []struct {
		name          string
		key           string
		label         string
		field         string
		opts          []func(*model.Filter) error
		expectedError error
	}{
		{
			name:          "basic terms filter",
			key:           "category",
			label:         "Category",
			field:         "category",
			opts:          nil,
			expectedError: nil,
		},
		{
			name:  "terms filter with default value",
			key:   "category",
			label: "Category",
			field: "category",
			opts: []func(*model.Filter) error{
				model.WithDefaultValue([]string{"electronics"}),
			},
			expectedError: nil,
		},
		{
			name:  "terms filter with value",
			key:   "category",
			label: "Category",
			field: "category",
			opts: []func(*model.Filter) error{
				model.WithValue([]any{"electronics", "books"}),
			},
			expectedError: nil,
		},
		{
			name:  "terms filter with non-string array elements",
			key:   "category",
			label: "Category",
			field: "category",
			opts: []func(*model.Filter) error{
				model.WithValue([]any{123, "books"}),
			},
			expectedError: nil,
		},
		{
			name:  "terms filter with empty array",
			key:   "category",
			label: "Category",
			field: "category",
			opts: []func(*model.Filter) error{
				model.WithValue([]any{}),
			},
			expectedError: nil,
		},
		{
			name:  "terms filter with nil array",
			key:   "category",
			label: "Category",
			field: "category",
			opts: []func(*model.Filter) error{
				model.WithValue(nil),
			},
			expectedError: nil,
		},
		// Negative test cases
		{
			name:  "terms filter with invalid value type",
			key:   "category",
			label: "Category",
			field: "category",
			opts: []func(*model.Filter) error{
				model.WithValue("invalid"),
			},
			expectedError: errorMsgForSetValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := model.NewFilter(dto.Filter{
				Key:           tt.key,
				Label:         tt.label,
				Field:         tt.field,
				OperationType: model.TermsOp,
			}, tt.opts...)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, errs.ErrInvalidArg))
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, filter)
			assert.Equal(t, tt.key, filter.Key())
			assert.Equal(t, tt.label, filter.Label())
			assert.Equal(t, tt.field, filter.Field())
			assert.Equal(t, model.TermsOp, filter.OperationType())
		})
	}
}

func TestRangeFilter(t *testing.T) {
	errorMsgForSetValue := fmt.Errorf("%w: value, filter key = 'price', operation type = 'range'", errs.ErrInvalidArg)
	tests := []struct {
		name          string
		key           string
		label         string
		field         string
		opts          []func(*model.Filter) error
		expectedError error
	}{
		{
			name:          "basic range filter",
			key:           "price",
			label:         "Price Range",
			field:         "price",
			opts:          nil,
			expectedError: nil,
		},
		{
			name:  "range filter with default value",
			key:   "price",
			label: "Price Range",
			field: "price",
			opts: []func(*model.Filter) error{
				model.WithDefaultValue(map[string]any{
					"from": 0.0,
					"to":   1000.0,
				}),
			},
			expectedError: nil,
		},
		{
			name:  "range filter with value",
			key:   "price",
			label: "Price Range",
			field: "price",
			opts: []func(*model.Filter) error{
				model.WithValue(map[string]any{
					"from": 10.0,
					"to":   100.0,
				}),
			},
			expectedError: nil,
		},
		// Negative test cases
		{
			name:  "range filter with invalid value type",
			key:   "price",
			label: "Price Range",
			field: "price",
			opts: []func(*model.Filter) error{
				model.WithValue("invalid"),
			},
			expectedError: errorMsgForSetValue,
		},
		{
			name:  "range filter with missing 'from' field",
			key:   "price",
			label: "Price Range",
			field: "price",
			opts: []func(*model.Filter) error{
				model.WithValue(map[string]any{
					"to": 100.0,
				}),
			},
			expectedError: nil,
		},
		{
			name:  "range filter with missing 'to' field",
			key:   "price",
			label: "Price Range",
			field: "price",
			opts: []func(*model.Filter) error{
				model.WithValue(map[string]any{
					"from": 10.0,
				}),
			},
			expectedError: nil,
		},
		{
			name:  "range filter with invalid 'from' type",
			key:   "price",
			label: "Price Range",
			field: "price",
			opts: []func(*model.Filter) error{
				model.WithValue(map[string]any{
					"from": "invalid",
					"to":   100.0,
				}),
			},
			expectedError: errorMsgForSetValue,
		},
		{
			name:  "range filter with invalid 'to' type",
			key:   "price",
			label: "Price Range",
			field: "price",
			opts: []func(*model.Filter) error{
				model.WithValue(map[string]any{
					"from": 10.0,
					"to":   "invalid",
				}),
			},
			expectedError: errorMsgForSetValue,
		},
		{
			name:  "range filter with 'from' greater than 'to'",
			key:   "price",
			label: "Price Range",
			field: "price",
			opts: []func(*model.Filter) error{
				model.WithValue(map[string]any{
					"from": 100.0,
					"to":   10.0,
				}),
			},
			expectedError: errorMsgForSetValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := model.NewFilter(dto.Filter{
				Key:           tt.key,
				Label:         tt.label,
				Field:         tt.field,
				OperationType: model.RangeOp,
			}, tt.opts...)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, errs.ErrInvalidArg))
				assert.Equal(t, tt.expectedError, err)

				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, filter)
			assert.Equal(t, tt.key, filter.Key())
			assert.Equal(t, tt.label, filter.Label())
			assert.Equal(t, tt.field, filter.Field())
			assert.Equal(t, model.RangeOp, filter.OperationType())
		})
	}
}

func TestMatchFilter(t *testing.T) {
	errorMsgForSetValue := fmt.Errorf("%w: value, filter key = 'name', operation type = 'match'", errs.ErrInvalidArg)

	tests := []struct {
		name          string
		key           string
		label         string
		field         string
		opts          []func(*model.Filter) error
		expectedError error
	}{
		{
			name:          "basic match filter",
			key:           "name",
			label:         "Name",
			field:         "name",
			opts:          nil,
			expectedError: nil,
		},
		{
			name:  "match filter with default value",
			key:   "name",
			label: "Name",
			field: "name",
			opts: []func(*model.Filter) error{
				model.WithDefaultValue("default search"),
			},
			expectedError: nil,
		},
		{
			name:  "match filter with value",
			key:   "name",
			label: "Name",
			field: "name",
			opts: []func(*model.Filter) error{
				model.WithValue("search value"),
			},
			expectedError: nil,
		},
		// Negative test cases
		{
			name:  "match filter with invalid value type (number)",
			key:   "name",
			label: "Name",
			field: "name",
			opts: []func(*model.Filter) error{
				model.WithValue(123),
			},
			expectedError: errorMsgForSetValue,
		},
		{
			name:  "match filter with invalid value type (boolean)",
			key:   "name",
			label: "Name",
			field: "name",
			opts: []func(*model.Filter) error{
				model.WithValue(true),
			},
			expectedError: errorMsgForSetValue,
		},
		{
			name:  "match filter with invalid value type (array)",
			key:   "name",
			label: "Name",
			field: "name",
			opts: []func(*model.Filter) error{
				model.WithValue([]string{"invalid"}),
			},
			expectedError: errorMsgForSetValue,
		},
		{
			name:  "match filter with invalid value type (map)",
			key:   "name",
			label: "Name",
			field: "name",
			opts: []func(*model.Filter) error{
				model.WithValue(map[string]any{"invalid": "value"}),
			},
			expectedError: errorMsgForSetValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := model.NewFilter(dto.Filter{
				Key:           tt.key,
				Label:         tt.label,
				Field:         tt.field,
				OperationType: model.MatchOp,
			}, tt.opts...)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, errs.ErrInvalidArg))
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, filter)
			assert.Equal(t, tt.key, filter.Key())
			assert.Equal(t, tt.label, filter.Label())
			assert.Equal(t, tt.field, filter.Field())
			assert.Equal(t, model.MatchOp, filter.OperationType())
		})
	}
}

func TestToggleFilter(t *testing.T) {
	errorMsgForSetValue := fmt.Errorf("%w: value, filter key = 'in_stock', operation type = 'toggle'", errs.ErrInvalidArg)

	tests := []struct {
		name          string
		key           string
		label         string
		field         string
		opts          []func(*model.Filter) error
		expectedError error
	}{
		{
			name:          "basic toggle filter",
			key:           "in_stock",
			label:         "In Stock",
			field:         "in_stock",
			opts:          nil,
			expectedError: nil,
		},
		{
			name:  "toggle filter with default value",
			key:   "in_stock",
			label: "In Stock",
			field: "in_stock",
			opts: []func(*model.Filter) error{
				model.WithDefaultValue(true),
			},
			expectedError: nil,
		},
		{
			name:  "toggle filter with value (bool)",
			key:   "in_stock",
			label: "In Stock",
			field: "in_stock",
			opts: []func(*model.Filter) error{
				model.WithValue(false),
			},
			expectedError: nil,
		},
		{
			name:  "toggle filter with value (float)",
			key:   "in_stock",
			label: "In Stock",
			field: "in_stock",
			opts: []func(*model.Filter) error{
				model.WithValue(1.0),
			},
			expectedError: nil,
		},
		{
			name:  "toggle filter with value (int)",
			key:   "in_stock",
			label: "In Stock",
			field: "in_stock",
			opts: []func(*model.Filter) error{
				model.WithValue(1),
			},
			expectedError: nil,
		},
		// Negative test cases
		{
			name:  "toggle filter with invalid value type (string)",
			key:   "in_stock",
			label: "In Stock",
			field: "in_stock",
			opts: []func(*model.Filter) error{
				model.WithValue("invalid"),
			},
			expectedError: errorMsgForSetValue,
		},
		{
			name:  "toggle filter with invalid value type (array)",
			key:   "in_stock",
			label: "In Stock",
			field: "in_stock",
			opts: []func(*model.Filter) error{
				model.WithValue([]string{"invalid"}),
			},
			expectedError: errorMsgForSetValue,
		},
		{
			name:  "toggle filter with invalid value type (map)",
			key:   "in_stock",
			label: "In Stock",
			field: "in_stock",
			opts: []func(*model.Filter) error{
				model.WithValue(map[string]any{"invalid": "value"}),
			},
			expectedError: errorMsgForSetValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := model.NewFilter(dto.Filter{
				Key:           tt.key,
				Label:         tt.label,
				Field:         tt.field,
				OperationType: model.ToggleOp,
			}, tt.opts...)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, errs.ErrInvalidArg))
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, filter)
			assert.Equal(t, tt.key, filter.Key())
			assert.Equal(t, tt.label, filter.Label())
			assert.Equal(t, tt.field, filter.Field())
			assert.Equal(t, model.ToggleOp, filter.OperationType())
		})
	}
}

func TestInvalidOperationType(t *testing.T) {
	filter, err := model.NewFilter(dto.Filter{
		Key:           "invalid",
		Label:         "Invalid",
		Field:         "invalid",
		OperationType: "invalid_op",
	})
	assert.Error(t, err)
	assert.Equal(t, "invalid argument: operation type = 'invalid_op', filter key = 'invalid'", err.Error())
	assert.Nil(t, filter)
}
