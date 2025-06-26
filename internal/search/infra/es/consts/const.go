package consts

const FAMILY_KEY_FIELD = "name"

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

const BUCKET_NAME_FOR_OFFSET = "offset"
const MAX_ALLOWED_FAMILY_QUANTITY = 1000
const MAX_ALLOWED_FAMILY_QUANTITY_FOR_OFFSET = 1000

var DefaultFamilyFields = []string{
	"id",
	"skus",
	"slug",
	"image",
	"countItems",
	"lowestAsk",
	"stockQty",
	"condition",
	"facilityAddressName",
	"productName",
	"productSlug",
	"productDisclaimer",
}

var SortingFields = []string{
	"lastSale",
	"newestListings",
	"expectedProfit",
	"expectedProfitPercent",
	"latestPurchased",
	"mostPurchased",
	"quantitySold",
	"recentlyReduced",
}
