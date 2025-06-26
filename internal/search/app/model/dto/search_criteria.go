package dto

type SearchCriteria struct {
	Filters    map[string]any
	OrderBy    *OrderBy
	Pagination *Pagination
}

type OrderBy struct {
	Field     string
	Direction string
}

type Pagination struct {
	Page uint64
	Size uint64
}
