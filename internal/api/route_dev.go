package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/api/dev"
)

func initDevRoutes(engine *gin.Engine) {

	group := engine.Group("")

	group.POST("/dev/op/init/users", wrapper(dev.OpCtr.InitUsers))

}
