/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package kit

import (
	"github.com/gin-gonic/gin"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao/mysql"
	"gorm.io/gorm"
)

// GetGormDB 获取GormDB
func GetGormDB(ctx *gin.Context) *gorm.DB {
	db, _ := ctx.Get(app.GormContext)
	return db.(*gorm.DB)
}

// GetNowUser 获取当前操作用户
func GetNowUser(ctx *gin.Context) *dbs.UmsUser {
	db := GetGormDB(ctx)
	usr := dbs.FindCache[dbs.UmsUser](db, ctx.GetUint64(app.GinUserID))
	return usr
}

// GetNowEmployee 获取当前操作员工，后台
func GetNowEmployee(ctx *gin.Context) (*dbs.UmsUser, *dbs.UmsEmployee) {
	db := GetGormDB(ctx)
	usr := dbs.FindCache[dbs.UmsUser](db, ctx.GetUint64(app.GinUserID))
	if usr == nil || usr.Type != app.UserTypeEmployee {
		return usr, nil
	}
	return usr, dbs.FindCache[dbs.UmsEmployee](db, usr.ID)
}
