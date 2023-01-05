/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package rds

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/mvity/go-boot/internal/conf"
)

var Redis *redis.Client

var RedisContext = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr:         conf.Config.Data.Redis.Addr,
		Username:     conf.Config.Data.Redis.Username,
		Password:     conf.Config.Data.Redis.Password,
		DB:           conf.Config.Data.Redis.Database,
		MinIdleConns: conf.Config.Data.Redis.MinIdle,
		MaxIdleConns: conf.Config.Data.Redis.MaxIdle,
	})
	if _, err := Redis.Ping(RedisContext).Result(); err != nil {
		return err
	}
	return nil
}
