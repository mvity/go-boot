package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-box/x"
	"github.com/mvity/go-quickstart/internal/app"
)

func InitApiService() error {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	return engine.Run(":" + x.ToString(app.Config.Port.ApiPort))
}
