/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package wss

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/wss/dev"
)

func initDevCmds(engine *gin.Engine) {

	engine.GET("/dev/ping", wrapper(dev.PingHandler.OnPingCommond))

}
