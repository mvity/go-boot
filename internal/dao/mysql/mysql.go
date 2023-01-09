/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package dbs

import (
	"context"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/conf"
	rds "github.com/mvity/go-boot/internal/dao/redis"
	"github.com/mvity/go-box/x"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var MySQL *gorm.DB

var MySQLContext = context.Background()

// InitMySQL 初始化MySQL组件
func InitMySQL() error {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,                                                // 慢 SQL 阈值
			LogLevel:                  x.Ternary(conf.Config.App.Debug, logger.Info, logger.Warn), // 日志级别
			IgnoreRecordNotFoundError: true,                                                       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,                                                       // 彩色打印
		},
	)

	mysqlConfig := mysql.Config{
		DSN:                       conf.Config.Data.MySQL.DSN,
		DefaultStringSize:         255,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}
	gromConfig := gorm.Config{
		Logger:                 newLogger,
		CreateBatchSize:        1024,
		AllowGlobalUpdate:      false, // 全局 update/delete
		SkipDefaultTransaction: true,  // 跳过默认事务
		NowFunc: func() time.Time {
			return time.Now().Local()
		}, // 更改创建时间使用的函数
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gromConfig); err != nil {
		return nil
	} else {
		MySQL = db
	}

	db, _ := MySQL.DB()
	db.SetMaxOpenConns(conf.Config.Data.MySQL.MaxOpen)                                         // 设置空闲连接池中连接的最大数量
	db.SetMaxIdleConns(conf.Config.Data.MySQL.MaxIdle)                                         // 设置打开数据库连接的最大数量。
	db.SetConnMaxLifetime(time.Duration(conf.Config.Data.MySQL.MaxConnLifetime) * time.Second) // 设置了连接可复用的最大时间。
	db.SetConnMaxIdleTime(time.Duration(conf.Config.Data.MySQL.MaxIdleTime) * time.Second)     // 设置空闲连接最大时间

	if err := MySQL.Callback().Create().Before("gorm:create").Register("BeforeCreateCallback", beforeCreateCallback); err != nil {
		return err
	}
	if err := MySQL.Callback().Update().Before("gorm:update").Register("BeforeUpdateCallback", beforeUpdateCallback); err != nil {
		return err
	}
	if err := MySQL.Callback().Update().After("gorm:update").Register("AfterUpdateCallback", afterUpdateCallback); err != nil {
		return err
	}
	if err := MySQL.Callback().Create().Before("*").Register("all_before", beforeCallback); err != nil {
		return err
	}
	if err := MySQL.Callback().Update().Before("*").Register("all_before", beforeCallback); err != nil {
		return err
	}
	if err := MySQL.Callback().Delete().Before("*").Register("all_before", beforeCallback); err != nil {
		return err
	}
	if err := MySQL.Callback().Query().Before("*").Register("all_before", beforeCallback); err != nil {
		return err
	}
	if err := MySQL.Callback().Row().Before("*").Register("all_before", beforeCallback); err != nil {
		return err
	}
	if err := MySQL.Callback().Raw().Before("*").Register("all_before", beforeCallback); err != nil {
		return err
	}
	return nil
}

/******* Callback start *******/

// 所有回调之前
func beforeCallback(db *gorm.DB) {
	if db.Statement.Context == nil || db.Statement.Context.Err() != nil {
		db.Statement.Context = MySQLContext
	}
}

// 创建之前
func beforeCreateCallback(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	ID := db.Statement.Schema.LookUpField("C001")
	if ID != nil {
		val, zero := ID.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		if zero || val == nil {
			db.Statement.SetColumn("C001", app.IDWorker.GetID())
		}
	}
	db.Statement.SetColumn("C002", 1)
	db.Statement.SetColumn("C003", false)
	db.Statement.SetColumn("C004", time.Now())
	db.Statement.SetColumn("C005", time.Now())

	var ctxMap map[string]string

	if ctxVal := db.Statement.Context.Value(app.GinContext); ctxVal != nil {
		ctxMap = ctxVal.(map[string]string)
	} else {
		ctxMap = make(map[string]string)
	}

	zyxCreateUid := db.Statement.Schema.LookUpField("C008")
	if zyxCreateUid != nil {
		val, zero := zyxCreateUid.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		if zero || val == nil {
			db.Statement.SetColumn("C008", ctxMap[app.GinUserID])
		}
	}
	zyxUpdateUid := db.Statement.Schema.LookUpField("C009")
	if zyxUpdateUid != nil {
		val, zero := zyxUpdateUid.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		if zero || val == nil {
			db.Statement.SetColumn("C009", ctxMap[app.GinUserID])
		}
	}
}

// 更新之前
func beforeUpdateCallback(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	zyxVersion := db.Statement.Schema.LookUpField("C002")
	if zyxVersion != nil {
		val, _ := zyxVersion.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		db.Statement.SetColumn("C002", val.(uint64)+1)
	}
	db.Statement.SetColumn("C005", time.Now())

	zyxUpdateUid := db.Statement.Schema.LookUpField("C009")
	if zyxUpdateUid != nil {
		if ctxVal := db.Statement.Context.Value(app.GinContext); ctxVal != nil {
			ctxMap := ctxVal.(map[string]string)
			db.Statement.SetColumn("C009", x.ToInt64(ctxMap[app.GinUserID]))
		} else {
			db.Statement.SetColumn("C009", app.PlatformID)
		}
	}

}

// 更新之后
func afterUpdateCallback(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	// 清空缓存
	ID := db.Statement.Schema.LookUpField("C001")
	if ID != nil {
		val, _ := ID.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		go rds.Cache.Clear("PK:"+db.Statement.Schema.Table, x.ToString(val))
	}
}

/******* Callback end  *******/
