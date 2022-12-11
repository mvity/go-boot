package rds

import (
	"github.com/mvity/go-boot/internal/dao/redis"
)

type count struct{}

// Count 计数器
var Count count

// Get 获取数量
func (c *count) Get(tag string) int64 {
	rkey := redis.DataPrefix + "Counter:" + tag
	return redis.Redis.IncrBy(redis.Context, rkey, 0).Val()
}

// Add  增加数量
func (c *count) Add(tag string, count int64) (bool, int64) {
	rkey := redis.DataPrefix + "Counter:" + tag
	val, err := redis.Redis.IncrBy(redis.Context, rkey, count).Result()
	return err == nil, val
}

// Sub  减少数量
func (c *count) Sub(tag string, count int64) (bool, int64) {
	rkey := redis.DataPrefix + "Counter:" + tag
	val, err := redis.Redis.DecrBy(redis.Context, rkey, count).Result()
	return err == nil, val
}

// Del 删除计数器
func (c *count) Del(tag string) {
	rkey := redis.DataPrefix + "Counter:" + tag
	redis.Redis.Del(redis.Context, rkey)
}
