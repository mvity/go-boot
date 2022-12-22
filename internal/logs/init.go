package logs

import (
	"github.com/mvity/go-boot/internal/conf"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
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
		Filename:   conf.Config.App.LogPath + "/" + tag + ".log",
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
