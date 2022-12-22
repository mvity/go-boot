package dao

import (
	"fmt"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-boot/internal/dao/dbs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMySQLDatabase 初始化数据库
func InitMySQLDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql",
		conf.Config.Data.MySQL.Username, conf.Config.Data.MySQL.Password,
		conf.Config.Data.MySQL.Host, conf.Config.Data.MySQL.Port)

	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
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

	if err := dbs.MySQL.AutoMigrate(&dbs.SysUser{}, &dbs.SysEmployee{}); err != nil {
		return err
	}

	return nil
}
