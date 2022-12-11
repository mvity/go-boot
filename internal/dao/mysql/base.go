package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/mvity/go-quickstart/internal/app"
	"github.com/mvity/go-quickstart/internal/dao/redis/rds"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
)

// 设置缓存
func setCache(id uint64, entity any) {
	if entity == nil {
		return
	}
	if ptr, ok := any(&entity).(DBEntity); ok && ptr.GetExpire().Milliseconds() == 0 {
		return
	}
	if bs, err := json.Marshal(entity); err == nil {
		if str := string(bs); len(str) > 0 {
			rds.Cache.Set("PK:"+(entity).(schema.Tabler).TableName(), strconv.FormatUint(id, 10), str)
		}
	}
}

// 清除缓存
func delCache(id uint64, entity any) {
	if ptr, ok := any(&entity).(DBEntity); ok && ptr.GetExpire().Milliseconds() == 0 {
		return
	}
	rds.Cache.Clear("PK:"+(entity).(schema.Tabler).TableName(), strconv.FormatUint(id, 10))
}

// 读取缓存
func getCache[T any](id uint64) (*T, bool) {
	if id <= 0 {
		return nil, false
	}
	var entity T
	if ptr, ok := any(&entity).(DBEntity); ok && ptr.GetExpire().Milliseconds() == 0 {
		return nil, false
	}
	if str := rds.Cache.Get("PK:"+any(&entity).(schema.Tabler).TableName(), strconv.FormatUint(id, 10)); len(str) > 0 {
		if err := json.Unmarshal([]byte(str), &entity); err == nil {
			return &entity, true
		} else {
			fmt.Printf("%v \n", err)
		}
	}
	return nil, false
}

// Save 保存
func Save(db *gorm.DB, entity any) {
	var idEntity *Entity
	if ptr, ok := entity.(DBEntity); ok {
		idEntity = ptr.GetEntity()
	}
	if idEntity == nil {
		panic(&app.MySQLError{Message: "无效的实体对象"})
	}
	if db.Statement.Context == nil || db.Statement.Context.Err() != nil {
		db.Statement.Context = MySQLContext
	}
	var result *gorm.DB
	if idEntity.ZyxVersion == 0 {
		result = db.Create(entity)
	} else {
		result = db.Model(entity).Where("C002 = ?", idEntity.ZyxVersion).Save(entity)
	}
	if result.Error != nil {
		panic(&app.MySQLError{Message: "保存失败", Origin: result.Error})
	}
	if result.RowsAffected == 0 {
		panic(&app.MySQLError{Message: "保存失败"})
	}
	delCache(idEntity.ID, entity)
}

// Remove 逻辑删除
func Remove(db *gorm.DB, entity any) {
	var idEntity *Entity
	if ptr, ok := entity.(DBEntity); ok {
		idEntity = ptr.GetEntity()
	}
	if idEntity == nil {
		panic(&app.MySQLError{Message: "无效的实体对象"})
	}
	if db.Statement.Context == nil || db.Statement.Context.Err() != nil {
		db.Statement.Context = MySQLContext
	}
	idEntity.ZyxDelete = true
	result := db.Model(entity).Where("C002 = ?", idEntity.ZyxVersion).Update("C003", true)
	if result.Error != nil {
		panic(&app.MySQLError{Message: "删除失败", Origin: result.Error})
	}
	if result.RowsAffected == 0 {
		panic(&app.MySQLError{Message: "删除失败"})
	}
	delCache(idEntity.ID, entity)
}

// Delete 物理删除
func Delete(db *gorm.DB, entity any) {
	var idEntity *Entity
	if ptr, ok := entity.(DBEntity); ok {
		idEntity = ptr.GetEntity()
	}
	if idEntity == nil {
		panic(&app.MySQLError{Message: "无效的实体对象"})
	}
	if db.Statement.Context == nil || db.Statement.Context.Err() != nil {
		db.Statement.Context = MySQLContext
	}
	result := db.Delete(entity)
	if result.Error != nil {
		panic(&app.MySQLError{Message: "删除失败", Origin: result.Error})
	}
	if result.RowsAffected == 0 {
		panic(&app.MySQLError{Message: "删除失败"})
	}
	delCache(idEntity.ID, entity)
}

// FindDatabase 获取数据库记录，包含删除标记的
func FindDatabase[T any](db *gorm.DB, id uint64) *T {
	if id <= 0 {
		return nil
	}
	var entity T
	if db.Statement.Context == nil || db.Statement.Context.Err() != nil {
		db.Statement.Context = MySQLContext
	}
	if db.Where("C001 = ?", id).Limit(1).Find(&entity).RowsAffected == 0 {
		return nil
	}
	var idEntity *Entity
	if ptr, ok := any(&entity).(DBEntity); ok {
		idEntity = ptr.GetEntity()
	}
	if idEntity == nil {
		panic(&app.MySQLError{Message: "无效的实体对象"})
	}
	if !idEntity.ZyxDelete {
		go setCache(idEntity.ID, &entity)
	}
	return &entity
}

// FindOrigin 查询数据库记录，不含删除标记的
func FindOrigin[T any](db *gorm.DB, id uint64) *T {
	if id <= 0 {
		return nil
	}
	var entity T
	if ptr := FindDatabase[T](db, id); ptr == nil {
		return nil
	} else {
		entity = *ptr
		var idEntity *Entity
		if obj, ok := any(&entity).(DBEntity); ok {
			idEntity = obj.GetEntity()
		}
		if idEntity == nil || idEntity.ZyxDelete {
			return nil
		}
		return &entity
	}
}

// FindCache 获取缓存数据，不存在查询未删除的数据库数据
func FindCache[T any](db *gorm.DB, id uint64) *T {
	if id <= 0 {
		return nil
	}
	var entity T
	if cache, ok := getCache[T](id); !ok {
		if ptr := FindOrigin[T](db, id); ptr == nil {
			return nil
		} else {
			entity = *ptr
		}
	} else {
		entity = *cache
	}
	return &entity
}

// FindSnapshot 获取缓存数据，不存在查询数据库数据
func FindSnapshot[T any](db *gorm.DB, id uint64) *T {
	if id <= 0 {
		return nil
	}
	var entity T
	if cache, ok := getCache[T](id); !ok {
		if ptr := FindDatabase[T](db, id); ptr == nil {
			return nil
		} else {
			entity = *ptr
		}
	} else {
		entity = *cache
	}
	return &entity
}

// findRecord 查询单条数据
func findRecord[T any](db *gorm.DB, query *Query) *T {
	var entity T
	if db.Raw(query.SQL, query.Param...).Limit(1).Scan(&entity).RowsAffected <= 0 {
		return nil
	}
	return &entity
}

// findRecords 查询多条数据
func findRecords[T any](db *gorm.DB, query *Query) []*T {
	var entitys = make([]*T, 0)
	db.Raw(query.SQL, query.Param...).Scan(&entitys)
	return entitys
}

// findPager 执行分页查询
//
//goland:noinspection ALL
func findPager[T any](db *gorm.DB, query *Query) (*app.Paged, []*T) {
	var countSQL, querySQL string
	cond := string([]rune(query.SQL)[strings.Index(query.SQL, " FROM ")+6:])
	{
		// 汇总
		countSQL = "SELECT " + "COUNT(C001)" + " FROM " + cond
		countSQL = strings.Split(countSQL, "ORDER BY")[0]
		var counter int
		db.Raw(countSQL, query.Param...).Scan(&counter)
		query.Count = counter
	}
	var entitys = make([]*T, 0)
	if query.Count > 0 && query.Count >= (query.Page-1)*query.Size {
		// 查询
		tmp := strings.Split(query.SQL, "WHERE")[0]
		limit := " LIMIT " + strconv.Itoa(query.Size*(query.Page-1)) + " , " + strconv.Itoa(query.Size)
		querySQL = tmp + "INNER JOIN (SELECT C001 FROM " + cond + limit + ") AS TMP USING(C001)" + query.Order
		db.Raw(querySQL, query.Param...).Scan(&entitys)
	}
	return query.Result(), entitys
}
