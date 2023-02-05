/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package dbs

import (
	"github.com/mvity/go-box/x"
	"time"
)

// DBEntity 实体接口
type DBEntity interface {
	// GetEntity 获取实体基础结构体
	GetEntity() *Entity
	// GetExpire 获取缓存失效时间，0不缓存
	GetExpire() time.Duration
}

// Entity 实体基结构体
type Entity struct {
	ID            uint64    `gorm:"column:C001;not null;primary_key;autoIncrement:false;comment:数据ID"`
	ZyxVersion    uint64    `gorm:"column:C002;not null;comment:数据版本"`
	ZyxDelete     bool      `gorm:"column:C003;not null;index;comment:删除标记"`
	ZyxCreateTime time.Time `gorm:"column:C004;not null;index;comment:创建时间"`
	ZyxUpdateTime time.Time `gorm:"column:C005;not null;index;comment:修改时间"`
}

// GetIDString 获取数据ID
func (e *Entity) GetIDString() string {
	return x.ToString(e.ID)
}

// Operator 操作人相关字段
type Operator struct {
	ZyxCreateUid uint64 `gorm:"column:C008;not null;index;comment:创建人ID"`
	ZyxUpdateUid uint64 `gorm:"column:C009;not null;index;comment:修改人ID"`
}
