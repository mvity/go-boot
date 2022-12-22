package app

import (
	"math"
	"strings"
	"time"
)

// Pager 分页查询参数
type Pager struct {
	Page int   `json:"page" binding:"gte=1" label:"分页页码"`
	Size int   `json:"size" binding:"gte=1,lte=10000" label:"加载数量"`
	Time int64 `json:"time" binding:"gte=0" label:"加载首页时间"`
}

// GenQuery 生成分页查询对象
func (p *Pager) GenQuery() *Query {
	return &Query{
		Page: p.Page,
		Size: p.Size,
		Time: p.Time,
	}
}

// Paged 分页查询结果
type Paged struct {
	Page  int   `json:"now"` // 当前分页
	Time  int64 `json:"fms"` // 初次查询时间
	Total int   `json:"all"` // 总页码
	Count int   `json:"row"` // 总记录数
}

// Query SQL查询对象
type Query struct {
	SQL   string // 业务SQL语句
	Param []any  // 业务SQL参数
	Order string // Order 语句

	Page int   // 当前分页
	Size int   // 分页数量
	Time int64 // 初次查询时间

	Total int // 总页码
	Count int // 总记录数
}

// AddSQL 添加语句
func (q *Query) AddSQL(sqlx string) *Query {
	q.SQL += " " + strings.TrimSpace(sqlx)
	return q
}

// AddOrder 添加语句
func (q *Query) AddOrder(sqlx string) *Query {
	q.Order = " " + strings.TrimSpace(sqlx)
	q.SQL += q.Order
	return q
}

// AddParam 添加参数
func (q *Query) AddParam(param any) *Query {
	q.Param = append(q.Param, param)
	return q
}

// AddSQLParam 添加语句和参数
func (q *Query) AddSQLParam(sqlx string, param ...any) *Query {
	q.AddSQL(sqlx)
	for _, a := range param {
		q.AddParam(a)
	}
	return q
}

// GenResult 生成分页结果
func (q *Query) GenResult() *Paged {
	var total = 1
	if q.Count > 0 {
		total = int(math.Ceil(float64(q.Count) / (float64(q.Size))))
	}
	if q.Page == 1 {
		q.Time = time.Now().UnixMilli()
	}
	return &Paged{
		Page:  q.Page,
		Time:  q.Time,
		Total: total,
		Count: q.Count,
	}
}
