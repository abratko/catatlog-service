package model

import "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/app/model/dto"

type SearchCriteria struct {
	filters    *Filters
	orderBy    OrderBy
	pagination Pagination
}

func WithPagination(pagination *dto.Pagination) func(*SearchCriteria) {
	return func(s *SearchCriteria) {

		if pagination == nil {
			return
		}

		if pagination.Page == 0 {
			return
		}

		if pagination.Size == 0 {
			return
		}

		s.pagination = Pagination{
			page:     pagination.Page,
			pageSize: pagination.Size,
		}
	}
}

func WithOrderBy(orderBy *dto.OrderBy) func(*SearchCriteria) {
	return func(s *SearchCriteria) {
		if orderBy == nil {
			return
		}

		if orderBy.Field == "" {
			return
		}

		direction := "asc"
		if orderBy.Direction == "desc" {
			direction = "desc"
		}

		s.orderBy = OrderBy{
			field:     orderBy.Field,
			direction: direction,
		}
	}
}

func WithFilters(filters *Filters) func(*SearchCriteria) {
	return func(s *SearchCriteria) {
		s.filters = filters
	}
}

func NewSearchCriteria(opts ...func(*SearchCriteria)) *SearchCriteria {
	s := &SearchCriteria{}
	s.pagination = Pagination{
		page:     1,
		pageSize: 40,
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s SearchCriteria) GetOrderBy() OrderBy {
	return s.orderBy
}

func (s SearchCriteria) GetPage() uint64 {
	return s.pagination.page
}

func (s SearchCriteria) GetPageSize() uint64 {
	return s.pagination.pageSize
}

func (s SearchCriteria) GetFilters() *Filters {
	return s.filters
}

type Pagination struct {
	page     uint64
	pageSize uint64
}

type OrderBy struct {
	field     string
	direction string
}

func (o OrderBy) GetField() string {
	return o.field
}

func (o OrderBy) GetDirection() string {
	return o.direction
}
