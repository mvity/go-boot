package wss

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-box/x"
)

// InitWssService 启动WebSocket服务
func InitWssService() error {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	initDevCmds(engine)

	logs.LogSysInfo("Start WebSocket service success, Port: "+x.ToString(conf.Config.Port.WebSocketPort), nil)
	return engine.Run(":" + x.ToString(conf.Config.Port.WebSocketPort))
}
