/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-box/x"
	"net/http"
	"strings"
)

// Handler 控制器方法
type Handler func(ctx *gin.Context) *app.Result

// wrapper 封装控制器方法
func wrapper(handler Handler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		result := handler(ctx)
		if result != nil {
			rString := x.JsonToString(result)
			if !ctx.GetBool(app.GinEncrypt) {
				ctx.JSON(http.StatusOK, result)
			} else {
				gReqId := ctx.Query("reqid")
				gTime := ctx.Query("time")
				aesKey := x.MD5String(gTime + gReqId)
				aesStart := x.ToInt(gTime[len(gTime)-1:])
				aesIv := x.MD5String(aesKey)[aesStart : aesStart+16]
				ctx.String(http.StatusOK, x.AESEncrypt(aesKey, aesIv, rString))
			}
			logs.LogApiInfo(ctx, result.Status.Error, rString)
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
