package agg

import (
	"fmt"
	"slices"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/infra/es/consts"
)

/**
 * Attaches sub aggregation to fetch family property
 */
func AppendSubAggAndSorting(
	criteria *model.SearchCriteria,
	originalAggs map[string]any,
) error {
	fieldsFromSortingRule := []string{}
	if criteria.GetOrderBy().GetField() != "" && slices.Contains(consts.SortingFields, criteria.GetOrderBy().GetField()) {
		fieldsFromSortingRule = []string{criteria.GetOrderBy().GetField()}
	}

	// we should attach fields from sorting
	// they should be part of the query to fetch families
	fields := append(consts.DefaultFamilyFields, fieldsFromSortingRule...)

	for _, fieldName := range fields {
		if fieldName == consts.FAMILY_KEY_FIELD || fieldName == consts.FAMILY_FIELDS["countItems"] {
			// Aggregation for name is on top level and
			// countItems is result of this aggregation
			// So, we don't need specials aggregation for these fields.
			continue
		}

		if _, ok := appenders[fieldName]; !ok {
			return fmt.Errorf("aggregation builder for field %s not found", fieldName)
		}

		appenders[fieldName](criteria, originalAggs)
	}

	return nil
}
