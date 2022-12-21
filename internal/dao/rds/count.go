package rds

import (
	"github.com/mvity/go-boot/internal/dao"
)

type count struct{}

// Count 计数器
var Count count

// Get 获取数量
func (c *count) Get(tag string) int64 {
	rkey := dao.RedisDataPrefix + "Counter:" + tag
	return dao.Redis.IncrBy(dao.MySQLContext, rkey, 0).Val()
}

// Add  增加数量
func (c *count) Add(tag string, count int64) (bool, int64) {
	rkey := dao.RedisDataPrefix + "Counter:" + tag
	val, err := dao.Redis.IncrBy(dao.MySQLContext, rkey, count).Result()
	return err == nil, val
}

// Sub  减少数量
func (c *count) Sub(tag string, count int64) (bool, int64) {
	rkey := dao.RedisDataPrefix + "Counter:" + tag
	val, err := dao.Redis.DecrBy(dao.MySQLContext, rkey, count).Result()
	return err == nil, val
}

// Del 删除计数器
func (c *count) Del(tag string) {
	rkey := dao.RedisDataPrefix + "Counter:" + tag
	dao.Redis.Del(dao.MySQLContext, rkey)
}
