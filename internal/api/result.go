package api

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/mvity/go-quickstart/internal/app"
	"time"
)

// Trans 翻译器
var Trans ut.Translator

// status 响应状态
type status struct {
	Error         int8   `json:"err"`
	Message       string `json:"msg,omitempty"`
	ServerNowTime string `json:"snt"`
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
			Error:         0,
			ServerNowTime: time.Now().Format("20060102150405"),
		},
	}
}

// Fail 响应失败
func Fail(error int8) *Result {
	return &Result{
		Status: &status{
			Error:         error,
			ServerNowTime: time.Now().Format("20060102150405"),
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

// SetPager 设置分页响应数据
func (r *Result) SetPager(paged *app.Paged, items any) *Result {
	if paged != nil && items != nil {
		data := make(map[string]any)
		data["pager"] = paged
		data["items"] = items
		r.Data = data
	}
	return r
}

// SetSummPager 设置分页和统计响应数据
func (r *Result) SetSummPager(paged *app.Paged, summ any, items any) *Result {
	if paged != nil && summ != nil && items != nil {
		data := make(map[string]any)
		data["pager"] = paged
		data["summ"] = summ
		data["items"] = items
		r.Data = data
	}
	return r
}
