package dto

type Item struct {
	Id         string
	Name       string
	Skus       []string
	Slug       string
	Image      string
	CountItems uint32
	LowestAsk  float64
	StockQty   uint32
	Conditions map[string]uint32
	Facilities map[string]uint32
	// This is info about any first a product
	ProductDisclaimer string
	ProductSlug       string
	ProductName       string
	SortingField      *SortingField
}

type SortingField struct {
	Field string
	Value string
}

type FacetValue struct {
	Value string
	Count uint32
}

type FacetedSearchResult struct {
	TotalCount uint32
	Items      []*Item
	Facets     map[string]FacetValue
}

type ItemsCollection struct {
	totalCount uint32
	items      []*Item
}

func NewFamilyTailCollection(totalCount uint32, items []*Item) *ItemsCollection {
	return &ItemsCollection{
		totalCount: totalCount,
		items:      items,
	}
}

func (i ItemsCollection) TotalCount() uint32 {
	return i.totalCount
}

func (i ItemsCollection) Items() []*Item {
	return i.items
}
