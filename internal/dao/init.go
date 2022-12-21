package dao

import (
	"fmt"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao/dbs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMySQLDatabase 初始化数据库
func InitMySQLDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql",
		app.Config.Data.MySQL.Username, app.Config.Data.MySQL.Password,
		app.Config.Data.MySQL.Host, app.Config.Data.MySQL.Port)

	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	} else {
		dbInit := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_0900_ai_ci'",
			app.Config.Data.MySQL.Database)
		db.Exec(dbInit)
	}
	return nil
}

// InitMySQLTable 初始化数据表
func InitMySQLTable() error {

	if err := dbs.MySQL.AutoMigrate(&dbs.SysUser{}); err != nil {
		return err
	}

	return nil
}
