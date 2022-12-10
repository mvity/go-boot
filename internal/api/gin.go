package api

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// Gin 错误定义
const (
	GinSignError        int8 = 10 // 签名校验错误
	GinTimeError        int8 = 11 // 客户端时间错误
	GinVersionError     int8 = 12 // 客户端版本过低
	GinServerClosed     int8 = 13 // 服务端暂停服务
	GinParamError       int8 = 14 // 参数校验失败
	GinAuthError        int8 = 15 // 身份鉴权失败
	GinActionError      int8 = 16 // 无接口操作权限
	GinDataError        int8 = 17 // 无数据操作权限
	GinNotFound         int8 = 18 // 未找到匹配接口
	GinMustParam        int8 = 19 // 缺少必要参数
	GinApiStop          int8 = 20 // 接口已停用
	GinApiLock          int8 = 21 // 请求业务处理中
	GinSysError         int8 = 22 // 系统处理异常
	GinTransactionError int8 = 99 // 数据处理异常
)

// Controller 控制器模型
type Controller func(ctx *gin.Context) *Result

// Index 根请求
func Index(*gin.Context) *Result {
	return Fail(GinApiStop).SetMessage("未开放此接口")
}

// NotFound 404
func NotFound(ctx *gin.Context) *Result {
	return Fail(GinNotFound).SetMessage("URI: [ " + strings.Split(ctx.Request.RequestURI, "?")[0] + " ] 无效")
}
