package rds

import (
	"github.com/mvity/go-boot/internal/dao"
	"time"
)

type locker struct{}

var Locker locker

// Lock 上锁
func (l *locker) Lock(tag string, second int) bool {
	rkey := dao.RedisLockPrefix + tag
	val, err := dao.Redis.SetNX(dao.MySQLContext, rkey, 1, time.Duration(second)*time.Second).Result()
	return err == nil && val
}

func (l *locker) UnLock(tag string) {
	rkey := dao.RedisLockPrefix + tag

	dao.Redis.Del(dao.MySQLContext, rkey)
}
