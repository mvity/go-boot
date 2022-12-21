package rds

import (
	"time"
)

type cache struct{}

// Cache 缓存器
var Cache cache

// Get 读取缓存数据
func (c *cache) Get(tag string, key string) string {
	rkey := RedisCachePrefix + tag + ":" + key
	return Redis.Get(RedisContext, rkey).Val()
}

// Set 设置缓存数据，有效期1天
func (c *cache) Set(tag string, key string, val string) {
	c.SetExpire(tag, key, val, 60*24)
}

// SetExpire 设置缓存数据，指定过期时间
func (c *cache) SetExpire(tag string, key string, val string, minute int) {
	rkey := RedisCachePrefix + tag + ":" + key
	Redis.SetEx(RedisContext, rkey, val, time.Duration(minute)*time.Minute)
}

// Clear 清除缓存
func (c *cache) Clear(tag string, key string) {
	rkey := RedisCachePrefix + tag + ":" + key
	go func() {
		for i := 0; i < 10; i++ {
			Redis.Del(RedisContext, rkey)
			time.Sleep(1 * time.Second)
		}
	}()
}
