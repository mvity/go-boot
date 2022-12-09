package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	k "github.com/mvity/go-box/kit"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"runtime"
	"time"
)

var (
	apiLogger, jobLogger, wssLogger, sysLogger, extLogger *zap.Logger
)

// InitLogger 初始化日志记录器
func InitLogger() error {
	handleLogger(apiLogger, "api")
	handleLogger(jobLogger, "job")
	handleLogger(wssLogger, "wss")
	handleLogger(sysLogger, "sys")
	handleLogger(extLogger, "ext")
	return nil
}

func handleLogger(logger *zap.Logger, tag string) {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	syncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   Config.App.LogPath + "/" + tag + ".log",
		MaxSize:    10,
		MaxAge:     180,
		MaxBackups: 180,
		Compress:   true,
		LocalTime:  true,
	})
	core := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)
	logger = zap.New(core, zap.AddCaller())
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			log.Fatalln(err)
		}
	}(logger)
}

// LogApiInfo 记录业务接口请求日志
func LogApiInfo(ctx *gin.Context, err int8, result string) {
	// 控制台输出
	if err > 0 {
		log.Printf("[FAIL] Api: URI: %v , Time: %v\nParam: \n%v\nResponse: \n%v\n", ctx.Request.URL, time.Since(ctx.GetTime(GinTime)), ctx.GetString(GinBody), result)
	} else if Config.App.Debug {
		log.Printf("[INFO] Api: URI: %v , Time: %v\nParam: \n%v\nResponse: \n%v\n", ctx.Request.URL, time.Since(ctx.GetTime(GinTime)), ctx.GetString(GinBody), result)
	}
	// 文件记录
	apiLogger.Info("Api invoke",
		zap.String("method", ctx.Request.Method),
		zap.String("url", ctx.Request.URL.String()),
		zap.String("header", k.ToJSONString(ctx.Request.Header)),
		zap.String("body", ctx.GetString(GinBody)),
		zap.Int8("err", err),
		zap.Uint64("uid", ctx.GetUint64(GinUserId)),
		zap.String("result", result),
		zap.Duration("dur", time.Since(ctx.GetTime(GinTime))),
	)
}

// LogExtInfo 记录请求第三方接口日志
func LogExtInfo(api string, uri string, param string, response string, status int, dur time.Duration) {
	// 控制台输出
	if Config.App.Debug {
		log.Printf("[INFO] [%s]: URI: %v, Param: %v, Status: %v, Response: %v, Time: %v", api, uri, param, status, response, dur)
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

// LogSysInfo 记录系统日志
func LogSysInfo(content string, err any) {
	if err == nil {
		pc, _, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)

		log.Printf("[INFO] Sys: <%v[%-4d]> %s", f.Name(), line, content)
		sysLogger.Info(fmt.Sprintf("<%v[%-4d]> %s", f.Name(), line, content))

	} else {
		max := 12
		if Config.App.Debug {
			max = 100
		}
		pc := make([]uintptr, max)
		n := runtime.Callers(0, pc)
		for i := 2; i < n; i++ {
			f := runtime.FuncForPC(pc[i])
			_, line := f.FileLine(pc[i])
			log.Printf("[FAIL] Sys: <%v[%-4d]> %v, %v", f.Name(), line, content, err)
			sysLogger.Error(fmt.Sprintf("<%v[%-4d]> %v, %v", f.Name(), line, content, err))
		}
	}
}

// LogNotifyInfo 记录回调接口日志
func LogNotifyInfo(biz string, param string, result any) {
	// 控制台输出
	if Config.App.Debug {
		log.Printf("[INFO] Notify [%s] , Param: %v , Result : %v", biz, param, result)
	}
	// 文件记录
	extLogger.Info("Notify invoke",
		zap.String("biz", biz),
		zap.String("param", param),
		zap.Any("result", result),
	)
}
