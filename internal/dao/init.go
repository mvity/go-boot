package dao

import "github.com/mvity/go-boot/internal/dao/dbs"

// InitMySQLEntity 初始化数据表
func InitMySQLEntity() error {

	if err := MySQL.AutoMigrate(&dbs.SysUser{}); err != nil {
		return err
	}

	return nil
}
