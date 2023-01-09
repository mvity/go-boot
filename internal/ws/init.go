/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-box/x"
)

var Server *WebsocketServer

// InitWssService 启动WebSocket服务
func InitWssService() error {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	Server = NewWsServer(conf.Config.Port.WebSocketPort)

	go Server.Start()

	engine.GET("/ws/:channel/:token", Server.Handler)

	logs.LogSysInfo("Start WebSocket service success, Port: "+x.ToString(conf.Config.Port.WebSocketPort), nil)
	return engine.Run(":" + x.ToString(conf.Config.Port.WebSocketPort))
}
