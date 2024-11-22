package util_test

import (
	"cake-scraper/pkg/util"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PaginatorSuite struct {
	suite.Suite
	data  []int
	slice func(offset, limit int64) []int
}

func (s *PaginatorSuite) TestPaginator() {
	// Given
	p := util.NewPaginator(s.slice, 1, 3, 10)

	// When
	items := p.Items()

	// Then
	s.Equal([]int{1, 2, 3}, items)

	p2 := util.NewPaginator(s.slice, 2, 3, 10)
	items2 := p2.Items()
	s.Equal([]int{4, 5, 6}, items2)
}

func (s *PaginatorSuite) TestPaginator_HasPrev() {
	// Given
	p := util.NewPaginator(s.slice, 1, 3, 10)

	// When
	hasPrev := p.HasPrev()

	// Then
	s.False(hasPrev)

	p2 := util.NewPaginator(s.slice, 2, 3, 10)
	hasPrev2 := p2.HasPrev()
	s.True(hasPrev2)
}

func (s *PaginatorSuite) TestPaginator_HasNext() {
	// Given
	slice := func(offset, limit int64) []int {
		return s.data[offset:limit]
	}
	p := util.NewPaginator(slice, 1, 3, 10)

	// When
	hasNext := p.HasNext()

	// Then
	s.True(hasNext)

	p2 := util.NewPaginator(slice, 4, 3, 10)
	hasNext2 := p2.HasNext()
	s.False(hasNext2)
}

func TestPaginatorSuite(t *testing.T) {
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i + 1
	}
	suite.Run(t, &PaginatorSuite{
		data: data,
		slice: func(offset, limit int64) []int {
			return data[offset : offset+limit]
		},
	})
}
