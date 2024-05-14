package paging

import (
	"math"
)

type PagingInput struct {
	Page         int `json:"page"`           // 当前页 (必须)
	PageSize     int `json:"page_size"`      // 每页数量  (没有 pageCount 时候必须存在此数据)
	DataTotal    int `json:"data_total"`     // 总数据量  (没有 pageCount 时候必须存在此数据)
	PrevBatch    int `json:"prev_batch"`     // 显示page前多少页
	NextBatch    int `json:"next_batch"`     // 显示page后多少页
	PrevSomePage int `json:"prev_some_page"` // 显示 page 前指定页
	NextSomePage int `json:"next_some_page"` // 显示 page 后指定页
}

type PagingOutput struct {
	HasPaging       bool  `json:"has_paging"`         // 是否存在分页
	PageCount       int   `json:"page_count"`         // 总页数
	DataTotal       int   `json:"data_total"`         // 总数据量
	Page            int   `json:"page"`               // 当前页
	IsFirstPage     bool  `json:"is_first_page"`      // 当前页是第一页
	IsLastPage      bool  `json:"is_last_page"`       // 当前页是最后一页
	PrevBatch       []int `json:"prev_batch"`         // 当前页前几页  不存在前几页则为 [] (根据 传入的 prevPages 扩展)
	NextBatch       []int `json:"next_batch"`         // 当前页后几页  不存在后几页则为 [] (根据 传入的 prevPages 扩展)
	PrevPage        int   `json:"prev_page"`          // 上一页
	NextPage        int   `json:"next_page"`          // 下一页
	PrevHasMorePage bool  `json:"prev_has_more_page"` // 除了 prevBatch 和 第一页还存在其他页
	NextHasMorePage bool  `json:"next_has_more_page"` // 除了 nextBatch 和 最后一页还存在其他页
	PrevSomePage    int   `json:"prev_some_page"`     // 当前页前 5 页 (根据 传入的 prevSomePage 决定是前几页)
	NextSomePage    int   `json:"next_some_page"`     // 当前页后 5 页 (根据 传入的 nextSomePage 决定是前几页)
	PageSize        int   `json:"page_size"`          // 每页显示数量
}

func Paginate(data PagingInput) (output PagingOutput) {
	// 默认值每页10条数据
	if data.PageSize < 1 {
		data.PageSize = 10
	}
	// 没有count和total无法计算分页 直接返回
	if data.DataTotal < 1 {
		return
	}
	output.PageCount = int(math.Ceil(float64(data.DataTotal) / float64(data.PageSize)))

	// page小于1
	if data.Page < 1 {
		data.Page = 1
	}
	// page大于pageCount
	if data.Page > output.PageCount {
		data.Page = output.PageCount
	}
	// 当前页
	output.Page = data.Page
	// 每页数量
	output.PageSize = data.PageSize
	// 总数量
	output.DataTotal = data.DataTotal
	// 存在分页
	output.HasPaging = true
	// page == 1 当前页为第一页 (page前面做过处理)
	if data.Page == 1 {
		output.IsFirstPage = true
	}
	// page == pageCount 为最后一页
	if data.Page == output.PageCount {
		output.IsLastPage = true
	}
	// 初始化数组
	output.PrevBatch = []int{}
	output.NextBatch = []int{}
	// 上一页  当前页不是第1页 并且总页数不是1
	if data.Page != 1 && output.PageCount != 1 {
		output.PrevPage = data.Page - 1
	} else {
		output.PrevPage = data.Page
	}
	if data.Page != output.PageCount {
		output.NextPage = data.Page + 1
	} else {
		output.NextPage = data.Page
	}
	// 当前页前x页 (没有值不进入)
	if data.PrevSomePage > 0 {
		output.PrevSomePage = data.Page - data.PrevSomePage
		// 防止 data.PrevSomePage > page
		if output.PrevSomePage <= 0 {
			output.PrevSomePage = 1
		}
	}
	// 当前页后x页
	if data.NextSomePage > 0 {
		output.NextSomePage = data.Page + data.NextSomePage
		// 防止 data.PrevSomePage > PageCount
		if output.NextSomePage > output.PageCount {
			output.NextSomePage = output.PageCount
		}
	}
	calculatePrevPage := data.Page - data.PrevBatch
	// 传入参数大于当前页
	//  2   2   [1]
	if calculatePrevPage <= 0 {
		calculatePrevPage = 1
	}
	// 最少3页才算有其他页
	if calculatePrevPage > 2 {
		output.PrevHasMorePage = true
	}
	// 首页必须显示 PrevBatch跳过首页
	if calculatePrevPage <= 1 {
		calculatePrevPage += 1
	}
	// 当前页前[]页
	if data.PrevBatch > 0 {
		// 1    1 []
		for i := calculatePrevPage; i < data.Page; i++ {
			output.PrevBatch = append(output.PrevBatch, i)
		}
	}
	// 当前页后[]页
	calculateNextPage := data.Page + data.NextBatch

	// 传入参数大于总页数
	if calculateNextPage > output.PageCount {
		calculateNextPage = output.PageCount
	}
	// 总页数 - 2 间隔1 有其他页
	if output.PageCount-2 >= calculateNextPage {
		output.NextHasMorePage = true
	}
	// 做过处理最大等于总页数 尾页必须显示 排除数组中
	if calculateNextPage == output.PageCount {
		calculateNextPage -= 1
	}
	if data.NextBatch > 0 {
		//             10    		5   总15
		// page + 1 跳过当前页
		for i := data.Page + 1; i <= calculateNextPage; i++ {
			output.NextBatch = append(output.NextBatch, i)
		}
	}
	return
}

func PaginateTwoVersions(data PagingInput) (output PagingOutput) {
	// 默认值每页10条数据
	if data.PageSize < 1 {
		data.PageSize = 10
	}
	// 没有count和total无法计算分页 直接返回
	if data.DataTotal < 1 {
		return
	}
	output.PageCount = int(math.Ceil(float64(data.DataTotal) / float64(data.PageSize)))

	// page小于1
	if data.Page < 1 {
		data.Page = 1
	}
	// page大于pageCount
	if data.Page > output.PageCount {
		data.Page = output.PageCount
	}
	// 当前页
	output.Page = data.Page
	// 每页数量
	output.PageSize = data.PageSize
	// 总数量
	output.DataTotal = data.DataTotal
	// 存在分页
	output.HasPaging = true
	// page == 1 当前页为第一页 (page前面做过处理)
	if data.Page == 1 {
		output.IsFirstPage = true
	}
	// page == pageCount 为最后一页
	if data.Page == output.PageCount {
		output.IsLastPage = true
	}
	// 初始化数组
	output.PrevBatch = []int{}
	output.NextBatch = []int{}
	// 上一页  当前页不是第1页 并且总页数不是1
	if data.Page != 1 && output.PageCount != 1 {
		output.PrevPage = data.Page - 1
	} else {
		output.PrevPage = data.Page
	}
	if data.Page != output.PageCount {
		output.NextPage = data.Page + 1
	} else {
		output.NextPage = data.Page
	}
	// 当前页前x页 (没有值不进入)
	if data.PrevSomePage > 0 {
		output.PrevSomePage = data.Page - data.PrevSomePage
		// 防止 data.PrevSomePage > page
		if output.PrevSomePage <= 0 {
			output.PrevSomePage = 1
		}
	}
	// 当前页后x页
	if data.NextSomePage > 0 {
		output.NextSomePage = data.Page + data.NextSomePage
		// 防止 data.PrevSomePage > PageCount
		if output.NextSomePage > output.PageCount {
			output.NextSomePage = output.PageCount
		}
	}
	if data.PrevBatch > 0 {
		// 自定义 首页
		firstPage := 1
		// 5  1   3  2
		calculateBatchIntervalValue := calculateBatchInterval(output.Page, firstPage)
		if calculateBatchIntervalValue < data.PrevBatch {
			calculateBatchIntervalValue = data.PrevBatch
		} else {
			calculateBatchIntervalValue = calculateBatchIntervalValue - data.PrevBatch
			output.PrevHasMorePage = true
		}
		for i := 1; i <= calculateBatchIntervalValue; i++ {
			output.PrevBatch = append(output.PrevBatch, output.Page)
		}
	}
	if data.NextBatch > 0 {
		// 自定义 首页
		firstPage := output.Page
		calculateBatchIntervalValue := calculateBatchInterval(output.PageCount, firstPage)
		if calculateBatchIntervalValue < data.NextBatch {
			calculateBatchIntervalValue = data.NextBatch
		} else {
			output.PrevHasMorePage = true
		}
		for i := 1; i <= calculateBatchIntervalValue; i++ {
			output.NextBatch = append(output.NextBatch, output.Page-i)
		}
	}

	return
}

func calculateBatchInterval(startPage, endPage int) (interval int) {
	var max, min int
	if startPage == endPage {
		return
	}
	if startPage > endPage {
		max, min = startPage, endPage
	}
	if startPage < endPage {
		max, min = endPage, startPage
	}
	return max - min - 1
}

// nolint:unused
func calculateBatch(maxScopeValue, batch int) (calculateBatchStartValue int) {
	// page => pageCount
	validBatchScopeMaxValue := maxScopeValue - 2
	// 当前页 2 - 2 有效 0
	if validBatchScopeMaxValue <= 0 {
		return
	}
	if validBatchScopeMaxValue <= batch {
		calculateBatchStartValue = validBatchScopeMaxValue
		return
	}
	// 总页数20  当前页 10
	// 当前页8   有效页6    前2页  = 6,7
	if validBatchScopeMaxValue > batch {
		calculateBatchStartValue = batch
		return
	}
	return
}
