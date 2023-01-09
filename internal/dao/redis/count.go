/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package rds

type count struct{}

// Count 计数器
var Count count

// Get 获取数量
func (c *count) Get(tag string) int64 {
	rkey := RedisDataPrefix + "Counter:" + tag
	return Redis.IncrBy(RedisContext, rkey, 0).Val()
}

// Add  增加数量
func (c *count) Add(tag string, count int64) (bool, int64) {
	rkey := RedisDataPrefix + "Counter:" + tag
	val, err := Redis.IncrBy(RedisContext, rkey, count).Result()
	return err == nil, val
}

// Sub  减少数量
func (c *count) Sub(tag string, count int64) (bool, int64) {
	rkey := RedisDataPrefix + "Counter:" + tag
	val, err := Redis.DecrBy(RedisContext, rkey, count).Result()
	return err == nil, val
}

// Del 删除计数器
func (c *count) Del(tag string) {
	rkey := RedisDataPrefix + "Counter:" + tag
	Redis.Del(RedisContext, rkey)
}
