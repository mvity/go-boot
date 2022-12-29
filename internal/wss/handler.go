/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package wss

import (
	"fmt"
	"github.com/mvity/go-box/x"
	"strings"
	"time"
)

// HandlerMessage 消息处理
func HandlerMessage(c *WsClient, message string) {
	c.DoHeartbeat(x.ToUInt64(time.Now().UnixMilli()))

	fmt.Printf("Receiver Data: %v , %v \n", c.Addr, string(message))
	if "ping" == strings.ToLower(message) {
		c.Send("pong")
	}

}
