/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package dbs

import (
	"github.com/mvity/go-boot/internal/app"
	"gorm.io/gorm"
	"strings"
	"time"
)

// UmsEmployee 平台用户信息表
type UmsEmployee struct {
	Entity
	Operator
	Roles    app.GormJSON `gorm:"column:B001;not null;type:json;comment:用户角色"`
	LoginKey string       `gorm:"column:B010;not null;index;size:32;comment:登录账号"`
	LoginPwd string       `gorm:"column:B011;not null;size:128;comment:登录密码"`
	Name     string       `gorm:"column:B020;not null;index;size:128;comment:用户名称"`
	Avatar   string       `gorm:"column:B021;not null;size:255;comment:用户头像地址"`
}

// TableName 数据表名称
func (e *UmsEmployee) TableName() string {
	return "UMS002"
}

// GetEntity 实体对象
func (e *UmsEmployee) GetEntity() *Entity {
	return &e.Entity
}

// GetExpire 缓存时长
func (e *UmsEmployee) GetExpire() time.Duration {
	return 24 * time.Hour
}

// UmsEmployeeFindByLoginKey 获取指定登录账号的平台用户
func UmsEmployeeFindByLoginKey(db *gorm.DB, loginKey string) *UmsEmployee {
	query := &app.Query{}
	query.AddSQL("SELECT C001 FROM UMS002 WHERE C003 = 0")
	if val := strings.TrimSpace(loginKey); val != "" {
		query.AddSQLParam("AND B010 = ?", val)
	} else {
		return nil
	}
	if ptr := findRecord[UmsEmployee](db, query); ptr != nil {
		return FindCache[UmsEmployee](db, ptr.ID)
	}
	return nil
}

// UmsEmployeeFindPager 查询分页数据
func UmsEmployeeFindPager(db *gorm.DB, query *app.Query, role string, loginKey string, name string) (*app.Pager, []*UmsEmployee) {
	query.AddSQL("SELECT C001 FROM UMS002 WHERE C003 = 0")
	if val := strings.TrimSpace(role); val != "" {
		query.AddSQLParam("AND B001 LIKE ?", "%"+val+"%")
	}
	if val := strings.TrimSpace(loginKey); val != "" {
		query.AddSQLParam("AND B010 = ?", val)
	}
	if val := strings.TrimSpace(name); val != "" {
		query.AddSQLParam("AND B020 = ?", val)
	}
	query.AddOrder("ORDER BY C001 DESC")

	pager, ids := findPager[UmsEmployee](db, query)
	var records []*UmsEmployee
	for _, rid := range ids {
		records = append(records, FindCache[UmsEmployee](db, rid.ID))
	}
	return pager, records
}
