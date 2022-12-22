package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/app"
)

type opCtr struct{}

var OpCtr opCtr

// InitUsers 初始化默认用户
func (*opCtr) InitUsers(ctx *gin.Context) *app.Result {

	return app.Succ()
}
