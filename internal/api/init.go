/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-box/x"
)

// InitApiService 启动Api服务
func InitApiService() error {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	engine.NoRoute(wrapper(NotFound))
	engine.NoMethod(wrapper(NotFound))
	/**
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
	*/
	engine.Use(cors.Default())
	engine.Use(errsHandler())
	engine.Use(initHandler())
	engine.Use(bodyHandler())
	engine.Use(signHandler())

	initValidator()

	initDevRoutes(engine)

	logs.LogSysInfo("Start API service success, Port: "+x.ToString(conf.Config.Port.ApiPort), nil)

	return engine.Run(":" + x.ToString(conf.Config.Port.ApiPort))
}
