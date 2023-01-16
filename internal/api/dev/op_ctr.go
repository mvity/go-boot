/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/app"
)

type initCtr struct{}

var InitCtr initCtr

// Users 初始化默认用户
func (*initCtr) Users(ctx *gin.Context) *app.Result {

	return app.Succ()
}
