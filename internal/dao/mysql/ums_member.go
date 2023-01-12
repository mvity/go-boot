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

// UmsMember 平台用户信息表
type UmsMember struct {
	Entity
	Operator
	Roles    app.GormJSON `gorm:"column:B001;not null;index;type:json;comment:用户角色"`
	LoginKey string       `gorm:"column:B010;not null;index;size:32;comment:登录账号"`
	LoginPwd string       `gorm:"column:B011;not null;size:128;comment:登录密码"`
	Mobile   string       `gorm:"column:B015;not null;index;size:16;comment:手机号"`
	Name     string       `gorm:"column:B020;not null;index;size:128;comment:用户名称"`
	Avatar   string       `gorm:"column:B021;not null;size:255;comment:用户头像地址"`
}

// TableName 数据表名称
func (e *UmsMember) TableName() string {
	return "UMS003"
}

// GetEntity 实体对象
func (e *UmsMember) GetEntity() *Entity {
	return &e.Entity
}

// GetExpire 缓存时长
func (e *UmsMember) GetExpire() time.Duration {
	return 24 * time.Hour
}

// UmsMemberFindByLoginKey 获取指定登录账号的平台用户
func UmsMemberFindByLoginKey(db *gorm.DB, loginKey string) *UmsMember {
	query := &app.Query{}
	query.AddSQL("SELECT C001 FROM UMS003 WHERE C003 = 0")
	if val := strings.TrimSpace(loginKey); val != "" {
		query.AddSQLParam("AND B010 = ?", val)
	} else {
		return nil
	}
	if ptr := findRecord[UmsMember](db, query); ptr != nil {
		return FindCache[UmsMember](db, ptr.ID)
	}
	return nil
}

// UmsMemberFindPager 查询分页数据
func UmsMemberFindPager(db *gorm.DB, query *app.Query, loginKey string, name string) (*app.Pager, []*UmsMember) {
	query.AddSQL("SELECT C001 FROM UMS003 WHERE C003 = 0")

	if val := strings.TrimSpace(loginKey); val != "" {
		query.AddSQLParam("AND B010 = ?", val)
	}
	if val := strings.TrimSpace(name); val != "" {
		query.AddSQLParam("AND B020 = ?", val)
	}
	query.AddOrder("ORDER BY C001 DESC")

	pager, ids := findPager[UmsMember](db, query)
	var records []*UmsMember
	for _, rid := range ids {
		records = append(records, FindCache[UmsMember](db, rid.ID))
	}
	return pager, records
}
