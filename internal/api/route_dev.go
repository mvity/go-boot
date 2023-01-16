/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/api/dev"
)

func initDevRoutes(engine *gin.Engine) {

	group := engine.Group("")

	group.POST("/dev/init/users", wrapper(dev.InitCtr.Users))

}
