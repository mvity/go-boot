package rds

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/mvity/go-boot/internal/app"
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
		Addr:         app.Config.Data.Redis.Host + ":" + x.ToString(app.Config.Data.Redis.Port),
		Password:     app.Config.Data.Redis.Password,
		DB:           x.ToInt(app.Config.Data.Redis.Database),
		PoolSize:     app.Config.Data.Redis.MaxConn,
		MaxIdleConns: app.Config.Data.Redis.MaxIdle,
		ReadTimeout:  time.Duration(app.Config.Data.Redis.Timeout),
	})
	if _, err := Redis.Ping(RedisContext).Result(); err != nil {
		return err
	}
	return nil
}
