/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package logs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-box/x"
	"go.uber.org/zap"
	"log"
	"runtime"
	"time"
)

// LogApiInfo 记录业务接口请求日志
func LogApiInfo(ctx *gin.Context, err int8, result string) {
	// 控制台输出
	if err > 0 {
		log.Printf("API [FAIL] URI: %v , ErrCode: %v , Duration: %v\nBody: %v\nResp: %v\n", ctx.Request.URL, err, time.Since(ctx.GetTime(app.GinTime)), x.StringDefaultIfBlank(ctx.GetString(app.GinBody), "{}"), result)
	} else if conf.Config.App.Debug {
		log.Printf("API [INFO] URI: %v , ErrCode: %v , Duration: %v\nBody: %v\nResp: %v\n", ctx.Request.URL, err, time.Since(ctx.GetTime(app.GinTime)), x.StringDefaultIfBlank(ctx.GetString(app.GinBody), "{}"), result)
	}
	// 文件记录
	apiLogger.Info("Api invoke",
		zap.String("method", ctx.Request.Method),
		zap.String("url", ctx.Request.URL.String()),
		zap.String("header", x.JsonToString(ctx.Request.Header)),
		zap.String("body", x.StringDefaultIfBlank(ctx.GetString(app.GinBody), "{}")),
		zap.Int8("err", err),
		zap.Uint64("uid", ctx.GetUint64(app.GinUserID)),
		zap.String("result", result),
		zap.Duration("dur", time.Since(ctx.GetTime(app.GinTime))),
	)
}

// LogJobInfo 记录任务日志
func LogJobInfo(content string, err any) {
	if err == nil {
		pc, _, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)

		log.Printf("JOB [INFO] <%v[%-4d]> %s", f.Name(), line, content)
		jobLogger.Info(fmt.Sprintf("<%v[%-4d]> %s", f.Name(), line, content))

	} else {
		max := 12
		if conf.Config.App.Debug {
			max = 100
		}
		pc := make([]uintptr, max)
		n := runtime.Callers(0, pc)
		for i := 2; i < n; i++ {
			f := runtime.FuncForPC(pc[i])
			_, line := f.FileLine(pc[i])
			log.Printf("JOB [FAIL] <%v[%-4d]> %v, %v", f.Name(), line, content, err)
			jobLogger.Error(fmt.Sprintf("<%v[%-4d]> %v, %v", f.Name(), line, content, err))
		}
	}
}

// LogWssInfo 记录Websocket日志
func LogWssInfo(addr string, userId uint64, message string) {
	// 控制台输出
	if conf.Config.App.Debug {
		log.Printf("WSS [INFO] [%s] , UserId: %v , %v", addr, userId, message)
	}
	// 文件记录
	wsLogger.Info("Wss invoke",
		zap.String("Addr", addr),
		zap.Uint64("UserId", userId),
		zap.String("Message", message),
	)
}

// LogSysInfo 记录系统日志
func LogSysInfo(content string, err any) {
	if err == nil {
		pc, _, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		log.Printf("SYS [INFO] <%v[%-4d]> %s", f.Name(), line, content)
		sysLogger.Info(fmt.Sprintf("<%v[%-4d]> %s", f.Name(), line, content))

	} else {
		max := 12
		if conf.Config.App.Debug {
			max = 100
		}
		pc := make([]uintptr, max)
		n := runtime.Callers(0, pc)
		for i := 2; i < n; i++ {
			f := runtime.FuncForPC(pc[i])
			_, line := f.FileLine(pc[i])
			log.Printf("SYS [FAIL] <%v[%-4d]> %v, %v", f.Name(), line, content, err)
			sysLogger.Error(fmt.Sprintf("<%v[%-4d]> %v, %v", f.Name(), line, content, err))
		}
	}
}

// LogExtInfo 记录请求第三方接口日志
func LogExtInfo(api string, uri string, param string, response string, status int, dur time.Duration) {
	// 控制台输出
	if conf.Config.App.Debug {
		log.Printf("EXT [INFO] [%s]: URI: %v, Param: %v, Status: %v, Response: %v, Time: %v", api, uri, param, status, response, dur)
	}
	// 文件记录
	extLogger.Info("Ext invoke",
		zap.String("api", api),
		zap.String("url", uri),
		zap.String("param", param),
		zap.String("resp", response),
		zap.Int("status", status),
		zap.Duration("dur", dur),
	)
}

// LogNotifyInfo 记录回调接口日志
func LogNotifyInfo(biz string, param string, result any) {
	// 控制台输出
	if conf.Config.App.Debug {
		log.Printf("Notify [INFO] [%s] , Param: %v , Result : %v", biz, param, result)
	}
	// 文件记录
	extLogger.Info("Notify invoke",
		zap.String("biz", biz),
		zap.String("param", param),
		zap.Any("result", result),
	)
}
