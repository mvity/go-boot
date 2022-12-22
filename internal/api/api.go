package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-box/x"
	"net/http"
	"strings"
)

// Controller 控制器模型
type Controller func(ctx *gin.Context) *app.Result

// 统一控制器模型
func wrapper(controller Controller) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		result := controller(ctx)
		if result != nil {
			noAES := ctx.GetBool(app.GinEncrypt) // 是否启用AES加密
			if noAES {
				ctx.JSON(http.StatusOK, result)
			} else {
				ctime := ctx.Query("time")
				reqid := ctx.Query("reqid")
				key := x.MD5String(ctime + reqid)
				iv := x.MD5String(key)[8:24]
				ctx.String(http.StatusOK, x.AESEncrypt(key, iv, x.JsonToString(result)))
			}
			logs.LogApiInfo(ctx, result.Status.Error, x.JsonToString(result))
		}
	}
}

// Index 根请求
func Index(*gin.Context) *app.Result {
	return app.Fail(app.GinApiStop).SetMessage("未开放此接口")
}

// NotFound 404
func NotFound(ctx *gin.Context) *app.Result {
	return app.Fail(app.GinNotFound).SetMessage("URI: [ " + strings.Split(ctx.Request.RequestURI, "?")[0] + " ] 无效")
}
