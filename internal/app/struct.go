package app

// Pager Gin分页参数
type Pager struct {
	Page int   `json:"page" binding:"gte=1" label:"分页页码"`
	Size int   `json:"size" binding:"gte=1,lte=10000" label:"加载数量"`
	Time int64 `json:"time" binding:"gte=0" label:"加载首页时间"`
}
