package rds

import (
	"github.com/mvity/go-boot/internal/dao/redis"
	"time"
)

type locker struct{}

var Locker locker

// Lock 上锁
func (l *locker) Lock(tag string, second int) bool {
	rkey := redis.LockPrefix + tag
	val, err := redis.Redis.SetNX(redis.Context, rkey, 1, time.Duration(second)*time.Second).Result()
	return err == nil && val
}

func (l *locker) UnLock(tag string) {
	rkey := redis.LockPrefix + tag

	redis.Redis.Del(redis.Context, rkey)
}
