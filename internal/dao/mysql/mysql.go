package mysql

import (
	"context"
	"fmt"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao/redis/rds"
	"github.com/mvity/go-box/x"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

var MySQL *gorm.DB

var MySQLContext = context.Background()

// InitMySQL 初始化MySQL组件
func InitMySQL() error {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,                                               // 慢 SQL 阈值
			LogLevel:                  x.Ternary(app.Config.App.Debug, logger.Info, logger.Warn), // 日志级别
			IgnoreRecordNotFoundError: true,                                                      // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,                                                      // 彩色打印
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		app.Config.Data.MySQL.Username, app.Config.Data.MySQL.Password,
		app.Config.Data.MySQL.Host, app.Config.Data.MySQL.Port, app.Config.Data.MySQL.Database)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         255,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}
	gromConfig := gorm.Config{
		Logger:                 newLogger,
		CreateBatchSize:        1024,
		AllowGlobalUpdate:      false,
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gromConfig); err != nil {
		return nil
	} else {
		MySQL = db
	}

	db, _ := MySQL.DB()
	db.SetMaxOpenConns(app.Config.Data.MySQL.MaxConn)
	db.SetMaxIdleConns(app.Config.Data.MySQL.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(app.Config.Data.MySQL.Timeout) * time.Minute)

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
			db.Statement.SetColumn("C008", ctxMap[app.GinUserId])
		}
	}
	zyxUpdateUid := db.Statement.Schema.LookUpField("C009")
	if zyxUpdateUid != nil {
		val, zero := zyxUpdateUid.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		if zero || val == nil {
			db.Statement.SetColumn("C009", ctxMap[app.GinUserId])
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
			db.Statement.SetColumn("C009", x.ToInt64(ctxMap[app.GinUserId]))
		} else {
			db.Statement.SetColumn("C009", app.PlatUserID)
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

// DBEntity 实体接口
type DBEntity interface {
	// GetEntity 获取实体基础结构体
	GetEntity() *Entity
	// GetExpire 获取缓存失效时间，0不缓存
	GetExpire() *time.Duration
}

// Entity 实体基结构体
type Entity struct {
	ID            uint64    `gorm:"column:C001;not null;primary_key;autoIncrement:false;comment:数据ID"`
	ZyxVersion    uint64    `gorm:"column:C002;not null;comment:数据版本"`
	ZyxDelete     bool      `gorm:"column:C003;not null;index;comment:删除标记"`
	ZyxCreateTime time.Time `gorm:"column:C004;not null;index;comment:创建时间"`
	ZyxUpdateTime time.Time `gorm:"column:C005;not null;index;comment:修改时间"`
}

// Operator 操作人相关字段
type Operator struct {
	ZyxCreateUid uint64 `gorm:"column:C008;not null;comment:创建人ID"`
	ZyxUpdateUid uint64 `gorm:"column:C009;not null;comment:修改人ID"`
}

// Query SQL查询对象
type Query struct {
	SQL   string // 业务SQL语句
	Param []any  // 业务SQL参数
	Order string // Order 语句

	Page int   // 当前分页
	Size int   // 分页数量
	Time int64 // 初次查询时间

	Total int // 总页码
	Count int // 总记录数
}

// AddSQL 添加语句
func (q *Query) AddSQL(sqlx string) *Query {
	q.SQL += " " + strings.TrimSpace(sqlx)
	return q
}

// AddOrder 添加语句
func (q *Query) AddOrder(sqlx string) *Query {
	q.Order = " " + strings.TrimSpace(sqlx)
	q.SQL += q.Order
	return q
}

// AddParam 添加参数
func (q *Query) AddParam(param any) *Query {
	q.Param = append(q.Param, param)
	return q
}

// AddSQLParam 添加语句和参数
func (q *Query) AddSQLParam(sqlx string, param ...any) *Query {
	q.AddSQL(sqlx)
	for _, a := range param {
		q.AddParam(a)
	}
	return q
}

// Result 生成分页结果
func (q *Query) Result() *app.Paged {
	var total = 1
	if q.Count > 0 {
		total = int(math.Ceil(float64(q.Count) / (float64(q.Size))))
	}
	if q.Page == 1 {
		q.Time = time.Now().UnixMilli()
	}
	return &app.Paged{
		Page:  q.Page,
		Time:  q.Time,
		Total: total,
		Count: q.Count,
	}
}
