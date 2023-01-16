/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package dao

import (
	"fmt"
	"github.com/mvity/go-boot/internal/conf"
	dbs "github.com/mvity/go-boot/internal/dao/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMySQLDatabase 初始化数据库
func InitMySQLDatabase() error {
	if db, err := gorm.Open(mysql.Open(conf.Config.Data.MySQL.DSN), &gorm.Config{}); err != nil {
		return err
	} else {
		dbInit := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_0900_ai_ci'",
			conf.Config.Data.MySQL.Database)
		db.Exec(dbInit)
	}
	return nil
}

// InitMySQLTable 初始化数据表
func InitMySQLTable() error {

	if err := dbs.MySQL.AutoMigrate(&dbs.UmsUser{}, &dbs.UmsEmployee{}); err != nil {
		return err
	}

	return nil
}
