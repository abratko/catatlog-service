package agg

import (
	"fmt"

	"gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model"
	stdAgg "gitlab.trgdev.com/gotrg/white-label/services/catalog/pkg/es/agg_factory"
)

var FAMILY_KEY_FIELD = "name"

var FAMILY_FIELDS = map[string]string{
	"id":                    "id",
	"skus":                  "skus",
	"slug":                  "slug",
	"image":                 "image",
	"countItems":            "countItems",
	"lowestAsk":             "lowestAsk",
	"lastSale":              "lastSale",
	"newestListings":        "newestListings",
	"expectedProfit":        "expectedProfit",
	"expectedProfitPercent": "expectedProfitPercent",
	"latestPurchased":       "latestPurchased",
	"mostPurchased":         "mostPurchased",
	"quantitySold":          "quantitySold",
	"recentlyReduced":       "recentlyReduced",
	"stockQty":              "stockQty",
	"condition":             "condition",
	"facilityAddressName":   "facilityAddressName",
	"productName":           "productName",
	"productSlug":           "productSlug",
	"productDisclaimer":     "productDisclaimer",
}

type RangeValue struct {
	From int
	To   int
}

type m = map[string]any

func appendParentFacilityAgg(
	criteria *model.SearchCriteria,
	aggs m,
) {
	filters := criteria.GetFilters()
	// the facility property inside product is nested document with self stock.qty property,
	// therefore aggregation by facility depends from inStock filter value.
	// If inStock filter is set, then we should show only facilities with stock.qty >= inStock.from
	minStock := 0
	if filters.Get("inStock").Value().(bool) {
		minStock = 1
	}

	facilityFilterInsideFacilityAgg := []m{
		{
			"range": m{
				"offers.stock.qty": m{
					"gte": minStock,
				},
			},
		},
	}

	location := filters.Get("location").Value()

	if location != nil {
		facilityFilterInsideFacilityAgg = append(
			facilityFilterInsideFacilityAgg,
			m{
				"terms": m{
					"offers.facility_address.id": location,
				},
			},
		)
	}

	aggs["aggs"].(m)["facility"] = m{
		"nested": m{
			"path": "offers",
		},
		"aggs": m{
			"facility": m{
				"filter": m{
					"bool": m{
						"filter": facilityFilterInsideFacilityAgg,
					},
				},
				"aggs": m{}, // here  will be child sub aggregations
			},
		},
	}
}

func appendFacilityNameAgg(
	criteria *model.SearchCriteria,
	rootAggs m,
) {
	rootFamilyAggs := rootAggs

	if _, ok := rootFamilyAggs["aggs"].(m)["facility"]; !ok {
		appendParentFacilityAgg(criteria, rootAggs)
	}

	rootFamilyAggs["aggs"].(m)["facility"].(m)["aggs"].(m)["facility"].(m)["aggs"].(m)[FAMILY_FIELDS["facilityAddressName"]] = stdAgg.TermsAgg(
		"offers.facility_address.name",
		stdAgg.WithSize(30),
		stdAgg.WithSubAggs(m{
			"stockQty": stdAgg.ValueAgg("offers.stock.qty", "sum"),
		}),
	)

	if criteria.GetOrderBy().GetField() != FAMILY_FIELDS["facilityAddressName"] {
		return
	}

	// count of facility address name and sorting by it
	rootFamilyAggs["aggs"].(m)["facility"].(m)["aggs"].(m)["facility"].(m)["aggs"].(m)["_count"] = stdAgg.ValueAgg("offers.facility_address.id", "cardinality")
	rootFamilyAggs["terms"].(m)["order"] = []m{
		{
			"facility>facility>_count.value": criteria.GetOrderBy().GetDirection(),
		},
	}
}

func appendConditionAgg(
	criteria *model.SearchCriteria,
	rootAggs m,
) {
	rootFamilyAggs := rootAggs
	rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["condition"]] = stdAgg.TermsAgg("condition")

	if criteria.GetOrderBy().GetField() != FAMILY_FIELDS["condition"] {
		return
	}

	rootFamilyAggs["aggs"].(m)["_conditionCount"] = stdAgg.ValueAgg("condition", "cardinality")

	rootFamilyAggs["terms"].(m)["order"] = []m{
		{"_conditionCount": criteria.GetOrderBy().GetDirection()},
	}
}

func appendSortingRule(
	field string,
	criteria *model.SearchCriteria,
	rootAggs m,
) {

	if criteria.GetOrderBy().GetField() != field {
		return
	}

	rootAggs["terms"].(m)["order"] = []m{
		{field: criteria.GetOrderBy().GetDirection()},
	}
}

func appendStockQtyAgg(
	criteria *model.SearchCriteria,
	rootAggs m,
) {

	if _, ok := rootAggs["aggs"].(m)["facility"]; !ok {
		appendParentFacilityAgg(criteria, rootAggs)
	}

	rootAggs["aggs"].(m)["facility"].(m)["aggs"].(m)["facility"].(m)["aggs"].(m)[FAMILY_FIELDS["stockQty"]] = stdAgg.ValueAgg("offers.stock.qty", "sum")

	if criteria.GetOrderBy().GetField() != FAMILY_FIELDS["stockQty"] {
		return
	}

	sortingFieldPath := fmt.Sprintf("facility>facility>%s.value", FAMILY_FIELDS["stockQty"])
	rootAggs["terms"].(m)["order"] = []m{
		{sortingFieldPath: criteria.GetOrderBy().GetDirection()},
	}
}

type AggAppender func(criteria *model.SearchCriteria, rootAggs m)

var appenders = map[string]AggAppender{
	FAMILY_FIELDS["condition"]:           appendConditionAgg,
	FAMILY_FIELDS["facilityAddressName"]: appendFacilityNameAgg,
	FAMILY_FIELDS["stockQty"]:            appendStockQtyAgg,
	FAMILY_FIELDS["productName"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootAggs["aggs"].(m)[FAMILY_FIELDS["productName"]] = stdAgg.TermsAgg(
			"name.raw",
			stdAgg.WithSize(1),
		)
	},
	FAMILY_FIELDS["productSlug"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootAggs["aggs"].(m)[FAMILY_FIELDS["productSlug"]] = stdAgg.TermsAgg("slug", stdAgg.WithSize(1))
	},
	FAMILY_FIELDS["productDisclaimer"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootAggs["aggs"].(m)[FAMILY_FIELDS["productDisclaimer"]] = stdAgg.TermsAgg("disclaimer", stdAgg.WithSize(1))

		appendSortingRule(FAMILY_FIELDS["productDisclaimer"], criteria, rootAggs)
	},
	FAMILY_FIELDS["id"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootAggs["aggs"].(m)[FAMILY_FIELDS["id"]] = stdAgg.TermsAgg("family.id", stdAgg.WithSize(1))

		appendSortingRule(FAMILY_FIELDS["id"], criteria, rootAggs)
	},
	FAMILY_FIELDS["skus"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootAggs["aggs"].(m)[FAMILY_FIELDS["skus"]] = stdAgg.TermsAgg("sku")

		appendSortingRule(FAMILY_FIELDS["skus"], criteria, rootAggs)
	},
	FAMILY_FIELDS["slug"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootAggs["aggs"].(m)[FAMILY_FIELDS["slug"]] = stdAgg.TermsAgg(
			"family.slug",
			stdAgg.WithSize(1),
		)
	},
	FAMILY_FIELDS["image"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootAggs["aggs"].(m)[FAMILY_FIELDS["image"]] = stdAgg.TermsAgg(
			"image",
			stdAgg.WithSize(1),
		)

	},
	FAMILY_FIELDS["lowestAsk"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootAggs["aggs"].(m)[FAMILY_FIELDS["lowestAsk"]] = stdAgg.ValueAgg("price", "min")

		appendSortingRule(FAMILY_FIELDS["lowestAsk"], criteria, rootAggs)
	},
	FAMILY_FIELDS["lastSale"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["lastSale"]] = stdAgg.ValueAgg("statistic.sell.total.last.price", "max")

		appendSortingRule(FAMILY_FIELDS["lastSale"], criteria, rootAggs)
	},
	FAMILY_FIELDS["newestListings"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["newestListings"]] = stdAgg.ValueAgg("sorting.newest_listings", "max")

		appendSortingRule(FAMILY_FIELDS["newestListings"], criteria, rootAggs)
	},
	FAMILY_FIELDS["expectedProfit"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["expectedProfit"]] = stdAgg.ValueAgg("sorting.expected_profit", "avg")

		appendSortingRule(FAMILY_FIELDS["expectedProfit"], criteria, rootAggs)
	},
	FAMILY_FIELDS["expectedProfitPercent"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["expectedProfitPercent"]] = stdAgg.ValueAgg("sorting.expected_profit_percent", "avg")

		appendSortingRule(FAMILY_FIELDS["expectedProfitPercent"], criteria, rootAggs)
	},
	FAMILY_FIELDS["latestPurchased"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["latestPurchased"]] = stdAgg.ValueAgg("sorting.latest_purchased", "max")

		appendSortingRule(FAMILY_FIELDS["latestPurchased"], criteria, rootAggs)
	},
	FAMILY_FIELDS["mostPurchased"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["mostPurchased"]] = stdAgg.ValueAgg("sorting.latest_purchased", "max")

		appendSortingRule(FAMILY_FIELDS["mostPurchased"], criteria, rootAggs)
	},
	FAMILY_FIELDS["quantitySold"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["quantitySold"]] = stdAgg.ValueAgg("sorting.quantity_sold", "sum")

		appendSortingRule(FAMILY_FIELDS["quantitySold"], criteria, rootAggs)
	},
	FAMILY_FIELDS["recentlyReduced"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["recentlyReduced"]] = stdAgg.ValueAgg("sorting.recently_reduced", "max")

		appendSortingRule(FAMILY_FIELDS["recentlyReduced"], criteria, rootAggs)
	},
	FAMILY_FIELDS["relevance"]: func(
		criteria *model.SearchCriteria,
		rootAggs m,
	) {
		rootFamilyAggs := rootAggs
		rootFamilyAggs["aggs"].(m)[FAMILY_FIELDS["relevance"]] = stdAgg.NewScriptAgg("_score", "max")

		appendSortingRule(FAMILY_FIELDS["relevance"], criteria, rootAggs)
	},
}
