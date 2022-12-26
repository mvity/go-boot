package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type pingHandler struct{}

var PingHandler pingHandler

func (*pingHandler) OnPingCommond(client *websocket.Conn, ctx *gin.Context) {
	err := client.WriteMessage(websocket.PingMessage, []byte("success"))
	if err != nil {
		return
	}
}
