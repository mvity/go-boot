package rds

import (
	"github.com/mvity/go-boot/internal/dao"
	"time"
)

type cache struct{}

// Cache 缓存器
var Cache cache

// Get 读取缓存数据
func (c *cache) Get(tag string, key string) string {
	rkey := dao.RedisCachePrefix + tag + ":" + key
	return dao.Redis.Get(dao.MySQLContext, rkey).Val()
}

// Set 设置缓存数据，有效期1天
func (c *cache) Set(tag string, key string, val string) {
	c.SetExpire(tag, key, val, 60*24)
}

// SetExpire 设置缓存数据，指定过期时间
func (c *cache) SetExpire(tag string, key string, val string, minute int) {
	rkey := dao.RedisCachePrefix + tag + ":" + key
	dao.Redis.SetEx(dao.MySQLContext, rkey, val, time.Duration(minute)*time.Minute)
}

// Clear 清除缓存
func (c *cache) Clear(tag string, key string) {
	rkey := dao.RedisCachePrefix + tag + ":" + key
	go func() {
		for i := 0; i < 10; i++ {
			dao.Redis.Del(dao.MySQLContext, rkey)
			time.Sleep(1 * time.Second)
		}
	}()
}
