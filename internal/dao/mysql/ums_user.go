/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package dbs

import (
	"time"
)

// UmsUser 用户信息
type UmsUser struct {
	Entity
	Operator
	Type            int8      `gorm:"column:B001;not null;index;comment:用户类型"`
	Name            string    `gorm:"column:B002;not null;size:128;comment:用户名称"`
	LastLoginTime   time.Time `gorm:"column:B010;not null;comment:最后登录时间"`
	LastLoginIP     string    `gorm:"column:B011;not null;size:64;comment:最后登录IP"`
	LastLoginCity   string    `gorm:"column:B012;not null;size:128;comment:最后登录城市"`
	LastLoginDevice string    `gorm:"column:B013;not null;size:128;comment:最后登录设备"`
	Lock            bool      `gorm:"column:B020;not null;index;comment:是否锁定登录"`
	LockTime        time.Time `gorm:"column:B021;not null;comment:锁定时间"`
	LockResumeTime  time.Time `gorm:"column:B022;not null;comment:解锁时间"`
	LockUserID      uint64    `gorm:"column:B023;not null;comment:锁定人ID"`
	LockCause       string    `gorm:"column:B024;not null;size:256;comment:锁定原因"`
}

// TableName 数据表名称
func (e *UmsUser) TableName() string {
	return "UMS001"
}

// GetEntity 实体对象
func (e *UmsUser) GetEntity() *Entity {
	return &e.Entity
}

// GetExpire 缓存时长
func (e *UmsUser) GetExpire() time.Duration {
	return 24 * time.Hour
}
