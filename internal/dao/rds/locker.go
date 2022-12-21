package rds

import (
	"time"
)

type locker struct{}

var Locker locker

// Lock 上锁
func (l *locker) Lock(tag string, second int) bool {
	rkey := RedisLockPrefix + tag
	val, err := Redis.SetNX(RedisContext, rkey, 1, time.Duration(second)*time.Second).Result()
	return err == nil && val
}

func (l *locker) UnLock(tag string) {
	rkey := RedisLockPrefix + tag

	Redis.Del(RedisContext, rkey)
}
