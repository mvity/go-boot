/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package wss

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-box/x"
	"net/http"
	"strings"
	"sync"
	"time"
)

// WsClient Websocket 客户端
type WsClient struct {
	ID             uint64          // 客户端ID
	Addr           string          // 客户端地址
	Connection     *websocket.Conn // 客户端连接会话
	Outbox         chan []byte     // 待发箱，等待发送的数据
	Channel        string          // 客户端连接通道，用于区分不同设备端连接
	UserId         uint64          // 客户端连接用户ID
	Auth           bool            // 是否鉴权成功
	ConnectionTime uint64          // 连接时间，首次连接时间
	HeartbeatTime  uint64          // 心跳时间，上次心跳时间
}

// NewWsClient 创建客户端
func NewWsClient(conn *websocket.Conn, channel string, userId uint64) *WsClient {
	now := x.ToUInt64(time.Now().UnixMilli())
	return &WsClient{
		ID:             app.IDWorker.GetID(),
		Addr:           conn.RemoteAddr().String(),
		Connection:     conn,
		Outbox:         make(chan []byte, 1024),
		Channel:        channel,
		UserId:         userId,
		Auth:           false,
		ConnectionTime: now,
		HeartbeatTime:  now,
	}
}

// GetClientKey 获取客户端标识
func (c *WsClient) GetClientKey() string {
	return fmt.Sprintf("%v_%v", strings.ToUpper(c.Channel), c.UserId)
}

// Read 读取客户端数据
func (c *WsClient) Read() {
	defer func() {
		Server.UnRegister <- c
		logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("Disconnect."))
		if err := recover(); err != nil {
			logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("ERR: Read error. %v", err))
		}
	}()
	for {

		if mt, message, err := c.Connection.ReadMessage(); err != nil {
			logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("ERR: %v", err))
			return
		} else {
			if mt == websocket.CloseMessage {
				break
			}
			c.Handle(string(message))
		}
	}
}

// Write 向客户端写入数据
func (c *WsClient) Write() {
	defer func() {
		logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("Disconnect."))
		if err := c.Connection.Close(); err != nil {
			logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("ERR: Disconnect error. %v", err))
		}
		if err := recover(); err != nil {
			logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("ERR: Write error. %v", err))
			if conf.Config.App.Debug {
				println(err)
			}
		}
	}()

	for {
		select {
		case message, ok := <-c.Outbox:
			if !ok {
				_ = c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
				logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("ERR: Write error. %v, %v", message, err))
				if conf.Config.App.Debug {
					println(err)
				}
				break
			}
		}
	}
}

// Close 关闭客户端连接
func (c *WsClient) Close() {
	err := c.Connection.Close()
	if err != nil {
		return
	}
	close(c.Outbox)
}

// Handle 处理收到消息
func (c *WsClient) Handle(message string) {
	c.DoHeartbeat(x.ToUInt64(time.Now().UnixMilli()))
	logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("Recv Data: %v", message))
	if "ping" == strings.ToLower(message) {
		c.Send("pong")
	}
}

// Send 发送信息
func (c *WsClient) Send(message string) {
	if c == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("ERR: SendMessage error. %v", err))
		}
	}()
	c.DoHeartbeat(x.ToUInt64(time.Now().UnixMilli()))
	c.Outbox <- []byte(message)
	logs.LogWssInfo(c.Addr, c.UserId, fmt.Sprintf("Send Data: %v", message))

}

// DoAuth 执行鉴权
func (c *WsClient) DoAuth() {
	c.Auth = true
}

// DoHeartbeat 执行心跳操作
func (c *WsClient) DoHeartbeat(now uint64) {
	c.HeartbeatTime = now
}

// CheckAlive 检查是否存活
func (c *WsClient) CheckAlive() bool {
	return c.HeartbeatTime+30*1000 >= x.ToUInt64(time.Now().UnixMilli())
}

// WsServer Websocket 服务端
type WsServer struct {
	Addr       string                 // 服务端地址
	Clients    map[string]*WsClient   // 客户端映射池
	Users      map[uint64][]*WsClient // 用户客户端关系
	Lock       sync.RWMutex           // 客户端映射池读写锁
	Register   chan *WsClient         // 客户端连接
	UnRegister chan *WsClient         // 客户端断开
	Outbox     chan []byte            // 待发箱，等待发送的广播数据
}

// NewWsServer 创建服务端
func NewWsServer(port int) *WsServer {
	return &WsServer{
		Addr:       "127.0.0.1:" + x.ToString(port),
		Clients:    make(map[string]*WsClient),
		Users:      make(map[uint64][]*WsClient),
		Register:   make(chan *WsClient, 1024),
		UnRegister: make(chan *WsClient, 1024),
		Outbox:     make(chan []byte, 1024),
	}
}

// AddClient 添加客户端
func (s *WsServer) AddClient(client *WsClient) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if old := s.Clients[client.GetClientKey()]; old != nil {
		old.Close()
	}
	s.Clients[client.GetClientKey()] = client
	ucs := s.Users[client.UserId]
	if ucs == nil {
		ucs = make([]*WsClient, 0)
	}
	ucs = append(ucs, client)
	s.Users[client.UserId] = ucs

}

// DelClient 移出客户端
func (s *WsServer) DelClient(client *WsClient) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Clients, client.GetClientKey())
	ucs := s.Users[client.UserId]
	_ucs := make([]*WsClient, 0)
	for _, uc := range ucs {
		if uc.ID != client.ID {
			_ucs = append(_ucs, uc)
		}
	}
	if len(_ucs) == 0 {
		delete(s.Users, client.UserId)
	} else {
		s.Users[client.UserId] = _ucs
	}
	client.Close()
}

// GetClient 获取指定通道和UserID的连接
func (s *WsServer) GetClient(channel string, userId uint64) *WsClient {
	return s.Clients[fmt.Sprintf("%v_%v", strings.ToUpper(channel), userId)]
}

// CheckClientAlive 检查客户端是否存活
func (s *WsServer) CheckClientAlive() {
	logs.LogWssInfo(s.Addr, app.PlatformID, "Current clients : "+x.ToString(len(s.Clients))+", "+x.ToString(len(s.Users)))
	for _, client := range s.Clients {
		if !client.CheckAlive() {
			s.DelClient(client)
		}
	}
	time.AfterFunc(30*time.Second, s.CheckClientAlive)
}

// Start 启动服务端
func (s *WsServer) Start() {
	go s.CheckClientAlive()
	for {
		select {
		case client := <-s.Register:
			s.AddClient(client)
		case client := <-s.UnRegister:
			s.DelClient(client)
		case message := <-s.Outbox:
			for _, client := range s.Clients {
				if !client.Auth {
					continue
				}
				select {
				case client.Outbox <- message:
				default:
					close(client.Outbox)
				}
			}
		}
	}
}

// Handler Gin处理函数  /ws/:channel/:uid
func (s *WsServer) Handler(ctx *gin.Context) {
	channel := ctx.Param("channel")
	uid := x.ToUInt64(ctx.Param("uid"))
	if channel == "" || uid == 0 {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	upgrader := websocket.Upgrader{
		HandshakeTimeout: 0,
		ReadBufferSize:   0,
		WriteBufferSize:  0,
		WriteBufferPool:  nil,
		Subprotocols:     []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
		Error:            nil,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}

	if conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		logs.LogWssInfo(ctx.Request.RemoteAddr, app.GuestID, fmt.Sprintf("ERR: Handler error. %v", err))
		http.NotFound(ctx.Writer, ctx.Request)
	} else {
		client := NewWsClient(conn, channel, uid)
		go client.Read()
		go client.Write()

		s.AddClient(client)
	}
}
