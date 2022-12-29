/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package logs

import (
	"github.com/mvity/go-boot/internal/conf"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var (
	apiLogger, jobLogger, wsLogger, sysLogger, extLogger *zap.Logger
)

// InitLogger 初始化日志记录器
func InitLogger() error {
	apiLogger = handleLogger("api")
	jobLogger = handleLogger("job")
	wsLogger = handleLogger("ws")
	sysLogger = handleLogger("sys")
	extLogger = handleLogger("ext")
	return nil
}

func handleLogger(tag string) *zap.Logger {

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
	logger := zap.New(core, zap.AddCaller())
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			log.Fatalln(err)
		}
	}(logger)
	return logger
}
