package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"time"
)

/*
 * Gin相关定义
 */

const (
	GinContext string = "Gin_Ctx"  // Gin Context 标识
	GinTime           = "Gin_Time" // Gin请求时间，服务端接收到请求时间
	GinBody           = "Gin_Body" // Gin请求内容，请求体解密后内容
	GinData           = "Gin_Data" // Gin请求内容，请求体原始内容，未解密
	GinLogger         = "Gin_Log"  // 是否记录请求日志
	GinEncrypt        = "Gin_AES"  // 是否进行AES加密
	GinUserId         = "Gin_Uid"  // 当前请求用户ID
)

// Controller 控制器模型
type Controller func(ctx *gin.Context) *Result

// Trans 翻译器
var Trans ut.Translator

// Pager Gin分页参数
type Pager struct {
	Page int   `json:"page" binding:"gte=1" label:"分页页码"`
	Size int   `json:"size" binding:"gte=1,lte=10000" label:"加载数量"`
	Time int64 `json:"time" binding:"gte=0" label:"加载首页时间"`
}

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
func (r *Result) SetPager(paged *Paged, items any) *Result {
	if paged != nil && items != nil {
		data := make(map[string]any)
		data["pager"] = paged
		data["items"] = items
		r.Data = data
	}
	return r
}

// SetSummPager 设置分页和统计响应数据
func (r *Result) SetSummPager(paged *Paged, summ any, items any) *Result {
	if paged != nil && summ != nil && items != nil {
		data := make(map[string]any)
		data["pager"] = paged
		data["summ"] = summ
		data["items"] = items
		r.Data = data
	}
	return r
}
