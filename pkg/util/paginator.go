package util

type Paginator[T any] interface {
	// Slice returns the slice function
	Slice(offset, limit int64) []T
	// CurrentPage returns the current page number
	CurrentPage() int64
	// Items returns the items in the current page
	Items() []T
	// HasPrev returns true if there is a previous page
	HasPrev() bool
	// Prev returns the previous page
	Prev() Paginator[T]
	// HasNext returns true if there is a next page
	HasNext() bool
	// Next returns the next page
	Next() Paginator[T]
	// Offset returns the offset of the current page
	Offset() int64
	// Count returns the number of items in the current page
	Count() int64
	// Total returns the total number of items
	Total() int64
	// ItemsPerPage returns the number of items per page
	PerPage() int64
	// TotalPage returns the total number of pages
	TotalPage() int64
}

type paginator[T any] struct {
	slice       func(offset, limit int64) []T
	currentPage int64
	perPage     int64
	total       int64
}

func NewPaginator[T any](slice func(offset, limit int64) []T, currentPage, perPage, total int64) Paginator[T] {
	return &paginator[T]{slice: slice, currentPage: currentPage, perPage: perPage, total: total}
}

func (p *paginator[T]) Slice(offset, limit int64) []T {
	offset = min(max(offset, 0), p.total)
	limit = max(min(limit, p.total-offset), 0)
	return p.slice(offset, limit)
}

func (p *paginator[T]) CurrentPage() int64 {
	return p.currentPage
}

func (p *paginator[T]) Items() []T {
	return p.slice(p.Offset(), p.Count())
}

func (p *paginator[T]) HasPrev() bool {
	return p.currentPage > 1
}

func (p *paginator[T]) Prev() Paginator[T] {
	return &paginator[T]{slice: p.slice, currentPage: p.currentPage - 1, total: p.total, perPage: p.perPage}
}

func (p *paginator[T]) HasNext() bool {
	return p.currentPage < p.TotalPage()
}

func (p *paginator[T]) Next() Paginator[T] {
	return &paginator[T]{slice: p.slice, currentPage: p.currentPage + 1, total: p.total, perPage: p.perPage}
}

func (p *paginator[T]) Offset() int64 {
	return (p.currentPage - 1) * p.perPage
}

func (p *paginator[T]) Count() int64 {
	return min(p.perPage, p.total-p.Offset())
}

func (p *paginator[T]) Total() int64 {
	return p.total
}

func (p *paginator[T]) PerPage() int64 {
	return p.perPage
}

func (p *paginator[T]) TotalPage() int64 {
	if p.total == 0 {
		return 0
	}
	return (p.total-1)/p.perPage + 1
}
