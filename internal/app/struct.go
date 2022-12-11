package app

// Pager 分页查询参数
type Pager struct {
	Page int   `json:"page" binding:"gte=1" label:"分页页码"`
	Size int   `json:"size" binding:"gte=1,lte=10000" label:"加载数量"`
	Time int64 `json:"time" binding:"gte=0" label:"加载首页时间"`
}

// Paged 分页查询结果
type Paged struct {
	Page  int   `json:"now"` // 当前分页
	Time  int64 `json:"fms"` // 初次查询时间
	Total int   `json:"all"` // 总页码
	Count int   `json:"row"` // 总记录数
}
