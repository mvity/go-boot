package wss

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/wss/dev"
)

func initDevCmds(engine *gin.Engine) {

	engine.GET("/dev/ping", wrapper(dev.PingHandler.OnPingCommond))

}
