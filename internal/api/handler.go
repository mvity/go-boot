/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-boot/internal/dao/dbs"
	"github.com/mvity/go-boot/internal/dao/rds"
	"github.com/mvity/go-boot/internal/kit"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-box/x"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 跨域处理
func corsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			ctx.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}

// 错误处理
func errsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var ae *app.ApiError
				if e1, ok := err.(*app.ApiError); ok {
					ae = e1
				} else if e2, ok := err.(*app.MySQLError); ok {
					ae = &app.ApiError{ErrCode: app.GinSysError, Message: "请求失败，数据库错误", Origin: e2}
				} else if e3, ok := err.(*app.RedisError); ok {
					ae = &app.ApiError{ErrCode: app.GinSysError, Message: "请求失败，缓存错误", Origin: e3}
				} else if e4, ok := err.(error); ok {
					ae = &app.ApiError{ErrCode: app.GinSysError, Message: "请求失败，系统错误", Origin: e4}
				} else {
					fmt.Println("Unknow error ", err)
					ae = &app.ApiError{ErrCode: app.GinSysError, Message: "请求失败，未知错误"}
				}
				result := app.Fail(ae.ErrCode).SetMessage(ae.Message)
				ctx.JSON(http.StatusOK, result)
				logs.LogApiInfo(ctx, result.Status.Error, x.JsonToString(result))
				if ae.Origin != nil {
					logs.LogSysInfo(ae.Message, ae.Origin)
				}
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

// 请求检查
func initHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "GET" {
			ctx.Next()
			return
		}
		if strings.Contains(ctx.Request.RequestURI, "/favicon.ico") {
			ctx.Next()
			return
		}

		ctx.Set(app.GinTime, time.Now())
		if strings.HasPrefix(ctx.Request.RequestURI, "/notify") {
			ctx.Set(app.GormContext, dbs.MySQL.WithContext(ctx.Request.Context()))
			ctx.Next()
			return
		}

		ctxMap := make(map[string]string)
		ctxMap[app.GinUserID] = x.ToString(app.GuestID)

		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), app.GinContext, ctxMap))
		ctx.Set(app.GormContext, dbs.MySQL.WithContext(ctx.Request.Context()))

		gToken := ctx.Query("token")
		gTime := ctx.Query("time")
		gNonce := ctx.Query("nonce")
		gSign := ctx.Query("sign")
		if gToken == "" || gTime == "" || gNonce == "" || gSign == "" {
			panic(&app.ApiError{ErrCode: app.GinMustParam, Message: "缺少必要参数"})
		}

		rtime, err := strconv.ParseInt(gTime, 10, 64)
		if err != nil {
			panic(&app.ApiError{ErrCode: app.GinTimeError, Message: "请求端时间错误", Origin: err})
		}
		if math.Abs(float64(rtime-time.Now().Unix())) > 5*60 {
			panic(&app.ApiError{ErrCode: app.GinTimeError, Message: "请求端时间错误"})
		}

		gVersion := ctx.Query("version")
		gVersionVal := x.ToInt(gVersion)
		if gVersionVal == 0 {
			panic(&app.ApiError{ErrCode: app.GinVersionError, Message: "请求端版本过低", Origin: err})
		}
		// 运营端请求
		if strings.HasPrefix(gVersion, "1") && gVersionVal < 0 {
			panic(&app.ApiError{ErrCode: app.GinVersionError, Message: "请求端版本过低", Origin: err})
		}

		if conf.Config.App.Debug {
			ctx.Set(app.GinEncrypt, ctx.DefaultQuery("aes", "true") == "true")
			ctx.Set(app.GinLogger, ctx.DefaultQuery("log", "true") == "true")
		}
		ctx.Next()
	}
}

// 请求Body 处理
func bodyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 非业务请求
		if ctx.Request.Method == "GET" {
			ctx.Next()
			return
		}
		if strings.HasPrefix(ctx.Request.RequestURI, "/notify") {
			ctx.Next()
			return
		}
		if strings.Contains(ctx.Request.RequestURI, "/upload") {
			ctx.Next()
			return
		}

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(&app.ApiError{ErrCode: app.GinParamError, Message: "参数解析失败", Origin: err})
		} else if len(body) == 0 {
			bodyNew := "{}"
			ctx.Set(app.GinBody, bodyNew)
			ctx.Set(app.GinData, bodyNew)
			bodyData := []byte(bodyNew)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyData))
			ctx.Request.ContentLength = int64(len(bodyData))
		} else {
			bodyStr := string(body)
			ctx.Set(app.GinData, bodyStr)
			if !strings.HasPrefix(bodyStr, "{") && !strings.HasPrefix(bodyStr, "[") {
				gNonce := ctx.Query("nonce")
				gTime := ctx.Query("time")
				aesKey := x.MD5String(gNonce + gTime)
				aesIv := x.MD5String(aesKey)[8:24]
				bodyStr = x.AESDecrypt(aesKey, aesIv, bodyStr)
			}
			var bodyNew = ""
			if strings.HasPrefix(bodyStr, "{") {
				var jmap map[string]any
				err = json.Unmarshal([]byte(bodyStr), &jmap)
				if err != nil {
					panic(&app.ApiError{ErrCode: app.GinParamError, Message: "参数解析失败", Origin: err})
				}
				for key, val := range jmap {
					if v, ok := val.(string); ok {
						jmap[key] = strings.TrimSpace(v)
					}
				}
				bodyNew = x.JsonToString(&jmap)
			} else {
				bodyNew = bodyStr
			}

			ctx.Set(app.GinBody, bodyNew)
			bodyData := []byte(bodyNew)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyData))
			ctx.Request.ContentLength = int64(len(bodyData))
		}
		ctx.Next()
	}
}

// 签名检查
func signHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 非业务请求
		if ctx.Request.Method == "GET" {
			ctx.Next()
			return
		}
		if strings.Contains(ctx.Request.RequestURI, "/favicon.ico") {
			ctx.Next()
			return
		}
		if strings.HasPrefix(ctx.Request.RequestURI, "/notify") {
			ctx.Next()
			return
		}
		if strings.Contains(ctx.Request.RequestURI, "/upload") {
			ctx.Next()
			return
		}

		gSign := ctx.Query("sign")
		if "312f74a079873e03a55c" == gSign {
			// debug 专用签名
			ctx.Next()
			return
		}
		gNonce := ctx.Query("nonce")
		gReqid := ctx.Query("reqid")
		gBody := ctx.GetString(app.GinData)
		if gSign != x.MD5String(gNonce+gBody+gReqid) {
			panic(&app.ApiError{ErrCode: app.GinSignError, Message: "请求签名错误"})
		}
		ctx.Next()
	}
}

// 设置NoAES标志
func noAesHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(app.GinEncrypt, false)
		ctx.Next()
	}
}

// // 设置NoLog标志
func noLogHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(app.GinLogger, false)
		ctx.Next()
	}
}

// 用户鉴权检查
func authHandler(auth bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.DefaultQuery("token", "Guest")
		uid := rds.Token.GetUserId(token)
		if uid < 0 {
			if auth {
				panic(&app.ApiError{ErrCode: app.GinAuthError, Message: "您的账号在其他地方登录"}) // 未登录用户
			}
		} else if uid == 0 {
			if auth {
				panic(&app.ApiError{ErrCode: app.GinAuthError, Message: "登录已失效，请重新登录"}) // 未登录用户
			}
		} else {
			if user := dbs.FindCache[dbs.SysUser](kit.GetGormDB(ctx), uint64(uid)); user != nil {
				ctxMap := (ctx.Request.Context().Value(app.GinContext)).(map[string]string)
				ctxMap[app.GinUserID] = strconv.FormatUint(user.ID, 10)
				ctx.Set(app.GinUserID, user.ID)
			} else {
				if auth {
					panic(&app.ApiError{ErrCode: app.GinAuthError, Message: "登录已失效，请重新登录"}) // 未登录用户
				}
			}

		}
		ctx.Next()
	}
}

// 用户类型检查
func typeHandler(types ...int8) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(types) > 0 {
			if user := dbs.FindCache[dbs.SysUser](kit.GetGormDB(ctx), ctx.GetUint64(app.GinUserID)); user == nil {
				panic(&app.ApiError{ErrCode: app.GinAuthError, Message: "登录已失效，请重新登录"}) // 未登录用户
			} else {
				flag := false
				for _, ut := range types {
					if user.Type == ut {
						flag = true
						break
					}
				}
				if !flag {
					panic(&app.ApiError{ErrCode: app.GinActionError, Message: "无操作权限"}) // 无接口权限
				}
			}
		}
		ctx.Next()
	}
}

// 用户角色检查
func roleHandler(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(roles) > 0 {
			if emp := dbs.FindCache[dbs.SysEmployee](kit.GetGormDB(ctx), ctx.GetUint64(app.GinUserID)); emp == nil {
				panic(&app.ApiError{ErrCode: app.GinAuthError, Message: "登录已失效，请重新登录"}) // 未登录用户
			} else {
				flag := false
				for _, role := range roles {
					if strings.Contains(emp.Roles.String(), role) {
						flag = true
						break
					}
				}
				if !flag {
					panic(&app.ApiError{ErrCode: app.GinActionError, Message: "无操作权限"}) // 无接口权限
				}
			}
		}
		ctx.Next()
	}
}
