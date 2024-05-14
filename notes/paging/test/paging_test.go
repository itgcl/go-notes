package paging__test

import (
	"testing"

	"go-notes/notes/paging"

	"github.com/stretchr/testify/assert"
)

// pageSize
func TestPageSize(t *testing.T) {
	{

		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  0,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageSize, 10)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  20,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageSize, 20)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  -10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageSize, 10)
	}
}

// page
func TestPage(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.Page, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:      0,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.Page, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:      -1,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.Page, 1)
	}
}

// pageCount
func TestPageCount(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageCount, 10)
	}
	{
		pagingData := paging.PagingInput{
			Page:      2,
			PageSize:  20,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageCount, 5)
	}
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageCount, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 99,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageCount, 10)
	}
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 101,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageCount, 11)
	}
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 0,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PageCount, 0)
	}
}

// DataTotal
func TestDataTotal(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 0,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.DataTotal, 0)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.DataTotal, 100)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: -1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.DataTotal, 0)
	}
}

// HasPaging
func TestHasPaging(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      0,
			PageSize:  0,
			DataTotal: 0,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.HasPaging, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.HasPaging, true)
	}
}

func TestIsFirstPage(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsFirstPage, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      10,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsFirstPage, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsFirstPage, true)
	}
	{
		pagingData := paging.PagingInput{
			Page:      -20,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsFirstPage, true)
	}
	{
		pagingData := paging.PagingInput{
			Page:      30,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsFirstPage, false)
	}
}

func TestIsLastPage(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsLastPage, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      9,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)

		assert.Equal(t, output.IsLastPage, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      10,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsLastPage, true)
	}
	{
		pagingData := paging.PagingInput{
			Page:      -20,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsLastPage, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      30,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.IsLastPage, true)
	}
}

func TestPrevPage(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevPage, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:      2,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)

		assert.Equal(t, output.PrevPage, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:      10,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevPage, 9)
	}
	{
		pagingData := paging.PagingInput{
			Page:      -1,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevPage, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  20,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevPage, 1)
	}
}

func TestNextPage(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextPage, 2)
	}
	{
		pagingData := paging.PagingInput{
			Page:      9,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)

		assert.Equal(t, output.NextPage, 10)
	}
	{
		pagingData := paging.PagingInput{
			Page:      10,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextPage, 10)
	}
	{
		pagingData := paging.PagingInput{
			Page:      20,
			PageSize:  10,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextPage, 10)
	}
	{
		pagingData := paging.PagingInput{
			Page:      10,
			PageSize:  20,
			DataTotal: 100,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextPage, 5)
	}
}

func TestPrevSomePage(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:         9,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 0,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 0)
	}
	{
		pagingData := paging.PagingInput{
			Page:         1,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:         1,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:         2,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:         3,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 2)
	}
	{
		pagingData := paging.PagingInput{
			Page:         8,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 3)
	}
	{
		pagingData := paging.PagingInput{
			Page:         20,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 5)
	}
	{
		pagingData := paging.PagingInput{
			Page:         3,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 4,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 1)
	}
	{
		pagingData := paging.PagingInput{
			Page:         5,
			PageSize:     10,
			DataTotal:    100,
			PrevSomePage: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevSomePage, 1)
	}
}

func TestNextSomePage(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:         9,
			PageSize:     10,
			DataTotal:    100,
			NextSomePage: 0,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextSomePage, 0)
	}
	{
		pagingData := paging.PagingInput{
			Page:         1,
			PageSize:     10,
			DataTotal:    100,
			NextSomePage: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextSomePage, 6)
	}
	{
		pagingData := paging.PagingInput{
			Page:         -1,
			PageSize:     10,
			DataTotal:    100,
			NextSomePage: 2,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextSomePage, 3)
	}
	{
		pagingData := paging.PagingInput{
			Page:         9,
			PageSize:     10,
			DataTotal:    100,
			NextSomePage: 2,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextSomePage, 10)
	}
}

func TestPrevBatch(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      9,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevBatch, []int{4, 5, 6, 7, 8})
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevBatch, []int{2, 3, 4})
	}
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevBatch, []int{})
	}
	{
		pagingData := paging.PagingInput{
			Page:      2,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevBatch, []int{})
	}
	{
		pagingData := paging.PagingInput{
			Page:      3,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevBatch, []int{2})
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  20,
			DataTotal: 100,
			PrevBatch: -1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevBatch, []int{})
	}
}

func TestNextBatch(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
			NextBatch: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextBatch, []int{2, 3, 4, 5, 6})
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
			NextBatch: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextBatch, []int{6, 7, 8, 9})
	}
	{
		pagingData := paging.PagingInput{
			Page:      10,
			PageSize:  10,
			DataTotal: 100,
			NextBatch: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextBatch, []int{})
	}
	{
		pagingData := paging.PagingInput{
			Page:      9,
			PageSize:  10,
			DataTotal: 100,
			NextBatch: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextBatch, []int{})
	}
	{
		pagingData := paging.PagingInput{
			Page:      8,
			PageSize:  10,
			DataTotal: 100,
			NextBatch: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextBatch, []int{9})
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  20,
			DataTotal: 100,
			NextBatch: -1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.NextBatch, []int{})
	}
}

func TestPrevHasMorePage(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:      8,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 5,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevHasMorePage, true)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 4,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevHasMorePage, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 3,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevHasMorePage, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 2,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevHasMorePage, true)
	}
	{
		pagingData := paging.PagingInput{
			Page:      10,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 10,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevHasMorePage, false)
	}
	{
		pagingData := paging.PagingInput{
			Page:      5,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: -1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevHasMorePage, true)
	}
	{
		pagingData := paging.PagingInput{
			Page:      1,
			PageSize:  10,
			DataTotal: 100,
			PrevBatch: 1,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.PrevHasMorePage, false)
	}
}

func TestOverall(t *testing.T) {
	{
		pagingData := paging.PagingInput{
			Page:         5,
			PageSize:     10,
			DataTotal:    100,
			PrevBatch:    4,
			NextBatch:    3,
			PrevSomePage: 1,
			NextSomePage: 2,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.Page, 5)
		assert.Equal(t, output.PageSize, 10)
		assert.Equal(t, output.DataTotal, 100)
		assert.Equal(t, output.PageCount, 10)
		assert.Equal(t, output.HasPaging, true)
		assert.Equal(t, output.IsFirstPage, false)
		assert.Equal(t, output.IsLastPage, false)
		assert.Equal(t, output.PrevPage, 4)
		assert.Equal(t, output.NextPage, 6)
		assert.Equal(t, output.PrevBatch, []int{2, 3, 4})
		assert.Equal(t, output.NextBatch, []int{6, 7, 8})
		assert.Equal(t, output.PrevSomePage, 4)
		assert.Equal(t, output.NextSomePage, 7)
		assert.Equal(t, output.PrevHasMorePage, false)
		assert.Equal(t, output.NextHasMorePage, true)
	}
	{
		pagingData := paging.PagingInput{
			Page:         1,
			PageSize:     20,
			DataTotal:    301,
			PrevBatch:    5,
			NextBatch:    6,
			PrevSomePage: 3,
			NextSomePage: 7,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.Page, 1)
		assert.Equal(t, output.PageSize, 20)
		assert.Equal(t, output.DataTotal, 301)
		assert.Equal(t, output.PageCount, 16)
		assert.Equal(t, output.HasPaging, true)
		assert.Equal(t, output.IsFirstPage, true)
		assert.Equal(t, output.IsLastPage, false)
		assert.Equal(t, output.PrevPage, 1)
		assert.Equal(t, output.NextPage, 2)
		assert.Equal(t, output.PrevBatch, []int{})
		assert.Equal(t, output.NextBatch, []int{2, 3, 4, 5, 6, 7})
		assert.Equal(t, output.PrevSomePage, 1)
		assert.Equal(t, output.NextSomePage, 8)
		assert.Equal(t, output.PrevHasMorePage, false)
		assert.Equal(t, output.NextHasMorePage, true)
	}
	{
		pagingData := paging.PagingInput{
			Page:         16,
			PageSize:     20,
			DataTotal:    305,
			PrevBatch:    8,
			NextBatch:    9,
			PrevSomePage: 15,
			NextSomePage: 2,
		}
		output := paging.Paginate(pagingData)
		assert.Equal(t, output.Page, 16)
		assert.Equal(t, output.PageSize, 20)
		assert.Equal(t, output.DataTotal, 305)
		assert.Equal(t, output.PageCount, 16)
		assert.Equal(t, output.HasPaging, true)
		assert.Equal(t, output.IsFirstPage, false)
		assert.Equal(t, output.IsLastPage, true)
		assert.Equal(t, output.PrevPage, 15)
		assert.Equal(t, output.NextPage, 16)
		assert.Equal(t, output.PrevBatch, []int{8, 9, 10, 11, 12, 13, 14, 15})
		assert.Equal(t, output.NextBatch, []int{})
		assert.Equal(t, output.PrevSomePage, 1)
		assert.Equal(t, output.NextSomePage, 16)
		assert.Equal(t, output.PrevHasMorePage, true)
		assert.Equal(t, output.NextHasMorePage, false)
	}
}
