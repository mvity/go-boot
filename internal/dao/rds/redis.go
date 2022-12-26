/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package rds

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-box/x"
	"time"
)

var Redis *redis.Client

var RedisContext = context.Background()

const (
	RedisBloomPrefix = "GoBoot:B:" // Redis 布隆过滤器前缀
	RedisCachePrefix = "GoBoot:C:" // Redis 缓存前缀
	RedisLockPrefix  = "GoBoot:L:" // Redis 锁前缀
	RedisDataPrefix  = "GoBoot:D:" // Redis 数据前缀
)

// InitRedis 初始化Redis连接
func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr:         conf.Config.Data.Redis.Host + ":" + x.ToString(conf.Config.Data.Redis.Port),
		Password:     conf.Config.Data.Redis.Password,
		DB:           x.ToInt(conf.Config.Data.Redis.Database),
		PoolSize:     conf.Config.Data.Redis.MaxConn,
		MaxIdleConns: conf.Config.Data.Redis.MaxIdle,
		ReadTimeout:  time.Duration(conf.Config.Data.Redis.Timeout),
	})
	if _, err := Redis.Ping(RedisContext).Result(); err != nil {
		return err
	}
	return nil
}
