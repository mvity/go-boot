/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package rds

import "github.com/mvity/go-boot/internal/conf"

var (
	RedisBloomPrefix = conf.Config.Data.Redis.Prefix + ":B:" // Redis 布隆过滤器前缀
	RedisCachePrefix = conf.Config.Data.Redis.Prefix + ":C:" // Redis 缓存前缀
	RedisLockPrefix  = conf.Config.Data.Redis.Prefix + ":L:" // Redis 锁前缀
	RedisDataPrefix  = conf.Config.Data.Redis.Prefix + ":D:" // Redis 数据前缀
)
