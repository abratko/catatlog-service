package model

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/errs"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"
)

type fetchOpts = func(ctx context.Context, optDtos any) (any, error)

type Filter struct {
	key           string
	label         string
	terms         *[]TermOpt
	range_        *RangeOpt
	field         string
	DefaultValue  any
	value         any
	operationType string
	isFaceted     bool
	fetchOpts     fetchOpts
}

func WithValue(value any) func(*Filter) error {
	return func(f *Filter) error {
		err := f.setValue(value)
		if err != nil {
			return err
		}
		return nil
	}
}

func WithDefaultValue(defaultValue any) func(*Filter) error {
	return func(f *Filter) error {
		// TODO: validate default value
		f.DefaultValue = defaultValue
		return nil
	}
}

func WithFetchOpts(fetchOpts fetchOpts) func(*Filter) error {
	return func(f *Filter) error {
		f.fetchOpts = fetchOpts
		return nil
	}
}

func NewFilter(
	dto dto.Filter,
	opts ...func(*Filter) error,
) (*Filter, error) {
	if !isValidOperationType(dto.OperationType) {
		return nil, fmt.Errorf("%w: operation type = '%s', filter key = '%s'", errs.ErrInvalidArg, dto.OperationType, dto.Key)
	}

	f := &Filter{
		key:           dto.Key,
		label:         dto.Label,
		field:         dto.Field,
		operationType: dto.OperationType,
		isFaceted:     dto.OperationType == TermsOp || dto.OperationType == RangeOp,
	}

	for _, opt := range opts {
		err := opt(f)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f Filter) Key() string {
	return f.key
}

func (f Filter) Label() string {
	return f.label
}

func (f Filter) Field() string {
	return f.field
}

func (f Filter) Terms() []TermOpt {
	return *f.terms
}

func (f Filter) Range() RangeOpt {
	return *f.range_
}

func (f Filter) OperationType() string {
	return f.operationType
}

func (f Filter) IsFaceted() bool {
	return f.isFaceted
}

func (f Filter) Value() any {
	if f.value == nil {
		return f.DefaultValue
	}

	return f.value
}

func (f *Filter) InitOpts(ctx context.Context, optDtos any) error {

	options := optDtos

	if f.fetchOpts != nil {
		var err error
		options, err = f.fetchOpts(ctx, optDtos)
		if err != nil {
			return err
		}
	}

	return f.setOptions(options)
}

func (f *Filter) setOptions(options any) error {
	errotMsg := fmt.Errorf("%w: options, filter key = '%s'", errs.ErrInvalidArg, f.Key())

	if options == nil {
		return nil
	}

	if f.OperationType() == TermsOp {
		terms, ok := options.([]TermOpt)
		if !ok {
			return errotMsg
		}
		f.terms = &terms

		return nil
	}

	if f.OperationType() == RangeOp {
		ranges, ok := options.(RangeOpt)
		if !ok {
			return errotMsg
		}
		f.range_ = &ranges

		return nil
	}

	return nil
}
func (f *Filter) setValue(value any) error {
	var ok bool
	errorMsg := fmt.Errorf("%w: value, filter key = '%s', operation type = '%s'", errs.ErrInvalidArg, f.Key(), f.OperationType())

	if isNil(value) {
		return nil
	}

	if f.OperationType() == TermsOp {
		if f.setValTerms(value) {
			return nil
		}

		return errorMsg
	}

	if f.OperationType() == RangeOp {
		if f.setRangeVal(value) {
			return nil
		}

		return errorMsg
	}

	if f.OperationType() == MatchOp {
		f.value, ok = value.(string)
		if ok {
			return nil
		}
		return errorMsg
	}

	switch v := value.(type) {
	case float64:
		f.value = v > 0
		return nil
	case bool:
		f.value = v
		return nil
	case int:
		f.value = v > 0
		return nil
	}

	return errorMsg
}

func (f *Filter) setRangeVal(value any) bool {
	newRange := RangeVal{}
	asMap, ok := value.(map[string]any)
	if !ok {
		return false
	}

	switch v := asMap["from"].(type) {
	case float64:
		newRange.From = &v
	case string:
		from, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
		newRange.From = &from
	case int:
		from := float64(v)
		newRange.From = &from
	}

	switch v := asMap["to"].(type) {
	case float64:
		newRange.To = &v
	case string:
		to, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
		newRange.To = &to
	case int:
		to := float64(v)
		newRange.To = &to
	}

	if newRange.From == nil && newRange.To == nil {
		return false
	}

	if newRange.From != nil && newRange.To != nil && *newRange.From > *newRange.To {
		return false
	}

	f.value = newRange

	return true
}

func (f *Filter) setValTerms(value any) bool {
	asArray, ok := value.([]any)
	if !ok {
		return false
	}

	values := make([]string, len(asArray))
	for i, v := range asArray {
		switch v := v.(type) {
		case string:
			values[i] = v
		case float64:
			values[i] = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			values[i] = strconv.Itoa(v)
		default:
			return false
		}
	}

	f.value = values
	return true
}

func (f *Filter) IsEmpty() bool {
	return isNil(f.value) && isNil(f.DefaultValue)
}

func (f *Filter) HasOptions() bool {
	if f.OperationType() == TermsOp {
		return f.terms != nil
	}

	if f.OperationType() == RangeOp {
		return f.range_ != nil
	}
	return false
}

type RangeOpt struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

type RangeVal struct {
	From *float64 `json:"from"`
	To   *float64 `json:"to"`
}

type TermOpt struct {
	Value string         `json:"value"`
	Count int            `json:"count"`
	Props map[string]any `json:"props"`
}

func isNil(value any) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case *int, *string, *float64, *bool:
		return v == nil
	default:
		return false
	}
}

func isValidOperationType(operationType string) bool {
	return operationType == TermsOp ||
		operationType == RangeOp ||
		operationType == MatchOp ||
		operationType == ToggleOp
}
