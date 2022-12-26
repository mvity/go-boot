package wss

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// Client Websocket 客户端
type Client struct {
}

// Server Websocket 服务端
type Server struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler Websocket处理器
type Handler func(client *websocket.Conn, ctx *gin.Context)

// wrapper 封装处理器方法
func wrapper(handler Handler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		client, _ := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		handler(client, ctx)
	}
}
