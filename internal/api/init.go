/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package api

import (
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

	engine.Use(corsHandler())
	engine.Use(errsHandler())
	engine.Use(initHandler())
	engine.Use(bodyHandler())
	engine.Use(signHandler())

	initValidator()

	initDevRoutes(engine)

	logs.LogSysInfo("Start API service success, Port: "+x.ToString(conf.Config.Port.ApiPort), nil)

	return engine.Run(":" + x.ToString(conf.Config.Port.ApiPort))
}
