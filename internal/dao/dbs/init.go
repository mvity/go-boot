package dbs

// InitMySQLEntity 初始化数据表
func InitMySQLEntity() error {

	if err := MySQL.AutoMigrate(&SysUser{}); err != nil {
		return err
	}

	return nil
}
