/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"math"
	"strings"
	"time"
)

// Pager 分页信息
type Pager struct {
	Page  int   `json:"page" binding:"gte=1" label:"分页页码"`
	Size  int   `json:"size" binding:"gte=1,lte=10000" label:"加载数量"`
	Time  int64 `json:"time" binding:"gte=0" label:"加载首页时间"`
	Pages int   `json:"pages" label:"总页码"`
	Rows  int   `json:"rows" label:"总记录数"`
}

// GenQuery 生成分页查询对象
func (p *Pager) GenQuery() *Query {
	return &Query{
		Page: p.Page,
		Size: p.Size,
		Time: p.Time,
	}
}

// Query SQL查询对象
type Query struct {
	SQL   string // 业务SQL语句
	Param []any  // 业务SQL参数
	Order string // Order 语句
	Page  int    // 当前分页
	Size  int    // 分页数量
	Time  int64  // 初次查询时间
	Pages int    // 总页码
	Rows  int    // 总记录数
}

// AddSQL 添加语句
func (q *Query) AddSQL(sql string) *Query {
	q.SQL += " " + strings.TrimSpace(sql)
	return q
}

// AddOrder 添加语句
func (q *Query) AddOrder(order string) *Query {
	q.Order = " " + strings.TrimSpace(order)
	q.SQL += q.Order
	return q
}

// AddParam 添加参数
func (q *Query) AddParam(param any) *Query {
	q.Param = append(q.Param, param)
	return q
}

// AddSQLParam 添加语句和参数
func (q *Query) AddSQLParam(sql string, param ...any) *Query {
	q.AddSQL(sql)
	for _, a := range param {
		q.AddParam(a)
	}
	return q
}

// GenPager 生成分页结果
func (q *Query) GenPager() *Pager {
	var pages = 1
	if q.Rows > 0 {
		pages = int(math.Ceil(float64(q.Rows) / (float64(q.Size))))
	}
	if q.Page == 1 {
		q.Time = time.Now().UnixMilli()
	}
	return &Pager{
		Page:  q.Page,
		Time:  q.Time,
		Pages: pages,
		Rows:  q.Rows,
	}
}

// Trans 翻译器
var Trans ut.Translator

// status 响应状态
type status struct {
	Error   int8   `json:"err"`
	Message string `json:"msg,omitempty"`
	Time    string `json:"snt"`
}

// Result 接口响应体
type Result struct {
	Status *status `json:"stat"`
	Data   any     `json:"data,omitempty"`
}

// Succ 响应成功
func Succ() *Result {
	return &Result{
		Status: &status{
			Error: 0,
			Time:  time.Now().Format("20060102150405"),
		},
	}
}

// Fail 响应失败
func Fail(error int8) *Result {
	return &Result{
		Status: &status{
			Error: error,
			Time:  time.Now().Format("20060102150405"),
		},
	}
}

// SetMessage 设置响应信息
func (r *Result) SetMessage(message string) *Result {
	if message != "" {
		r.Status.Message = message
	}
	return r
}

// SetError 设置Error错误
func (r *Result) SetError(err error) *Result {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var maps = errs.Translate(Trans)
		var msgs string
		for _, v := range maps {
			msgs += v + "；"
		}
		r.Status.Message = fmt.Sprintf("%v", msgs)
	} else {
		r.Status.Message = fmt.Sprintf("%v", err)
	}
	return r
}

// SetAttr 设置响应数据，单个字段
func (r *Result) SetAttr(key string, value any) *Result {
	if value != nil && key != "" {
		data := make(map[string]any)
		data[key] = value
		r.Data = data
	}
	return r
}

// SetItem 设置响应数据，单条
func (r *Result) SetItem(item any) *Result {
	if item != nil {
		r.Data = item
	}
	return r
}

// SetItems 设置响应数据，多条
func (r *Result) SetItems(items any) *Result {
	if items != nil {
		item := make(map[string]any)
		item["items"] = items
		r.Data = item
	}
	return r
}

// SetPagerItems 设置分页响应数据
func (r *Result) SetPagerItems(pager *Pager, items any) *Result {
	if pager != nil && items != nil {
		data := make(map[string]any)
		data["pager"] = pager
		data["items"] = items
		r.Data = data
	}
	return r
}

// SetSummPagerItems 设置分页和统计响应数据
func (r *Result) SetSummPagerItems(pager *Pager, summ any, items any) *Result {
	if pager != nil && summ != nil && items != nil {
		data := make(map[string]any)
		data["pager"] = pager
		data["summ"] = summ
		data["items"] = items
		r.Data = data
	}
	return r
}
