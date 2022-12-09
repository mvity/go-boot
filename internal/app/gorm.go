package app

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"math"
	"strings"
	"time"
)

/*
 * Gorm 相关定义
 */

const (
	GormContext string = "Gorm_Ctx" // Gorm Context 标识
)

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

// Paged 分页查询结果
type Paged struct {
	Page  int   `json:"now"` // 当前分页
	Time  int64 `json:"fms"` // 初次查询时间
	Total int   `json:"all"` // 总页码
	Count int   `json:"row"` // 总记录数
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
func (q *Query) Result() *Paged {
	var total = 1
	if q.Count > 0 {
		total = int(math.Ceil(float64(q.Count) / (float64(q.Size))))
	}
	if q.Page == 1 {
		q.Time = time.Now().UnixMilli()
	}
	return &Paged{
		Page:  q.Page,
		Time:  q.Time,
		Total: total,
		Count: q.Count,
	}
}

// GormJSON defined JSON data type, need to implement driver.Valuer, sql.Scanner interface
type GormJSON json.RawMessage

// Value return json value, implement driver.Valuer interface
func (j GormJSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	bytes, err := json.RawMessage(j).MarshalJSON()
	return string(bytes), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *GormJSON) Scan(value any) error {
	if value == nil {
		*j = GormJSON("null")
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = GormJSON(result)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (j GormJSON) MarshalJSON() ([]byte, error) {
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON to deserialize []byte
func (j *GormJSON) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := result.UnmarshalJSON(b)
	*j = GormJSON(result)
	return err
}

func (j GormJSON) String() string {
	return string(j)
}

// GormDataType gorm common data type
func (j GormJSON) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (j GormJSON) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (j GormJSON) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	if len(j) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := j.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if _, ok := db.Dialector.(*mysql.Dialector); ok {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}
