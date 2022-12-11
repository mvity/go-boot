package rds

import (
	"errors"
	"fmt"
	"github.com/mvity/go-box/x"
	"github.com/mvity/go-quickstart/internal/dao/redis"
	"strconv"
	"strings"
	"time"
)

type serial struct{}

var Serial serial

// Next 获取下一个序列号
func (*serial) Next(tag string, prifix string, size int, expireMinutes int64) string {
	rKey := redis.RedisDataPrefix + "Serial:Seq:" + tag + ":" + prifix
	max, _ := strconv.ParseInt(strings.Repeat("9", size), 10, 64)
	init := redis.Redis.Exists(redis.RedisContext, rKey).Val() > 0
	val := redis.Redis.IncrBy(redis.RedisContext, rKey, 1).Val()
	if val > max {
		panic(errors.New("out of maximum serial number"))
	}
	if !init {
		redis.Redis.Expire(redis.RedisContext, rKey, time.Duration(expireMinutes)*time.Minute)
	}
	return fmt.Sprintf("%0"+strconv.FormatInt(int64(size), 10)+"d", val)
}

// Random 获取下一个随机序列号
func (*serial) Random(tag string, prifix string, size int, expireMinutes int64) string {
	rKey := redis.RedisDataPrefix + "Serial:Rdm:" + tag + ":" + prifix
	init := redis.Redis.Exists(redis.RedisContext, rKey).Val() > 0
	val := x.RandomString(size, false, true)
	for i := 0; i < size; i++ {
		if redis.Redis.SIsMember(redis.RedisContext, rKey, val).Val() {
			val = x.RandomString(size, false, true)
		} else {
			redis.Redis.SAdd(redis.RedisContext, rKey, val)
			break
		}
	}
	if !init {
		redis.Redis.Expire(redis.RedisContext, rKey, time.Duration(expireMinutes)*time.Minute)
	}
	return val
}

// RandomFixed 获取下一个随机序列号
func (*serial) RandomFixed(tag string, min int64, max int64) int64 {
	rKey := redis.RedisDataPrefix + "Serial:Rdm:Fixed:" + tag
	val := x.RandomInt(min, max)
	for i := min; i < max; i++ {
		if redis.Redis.SIsMember(redis.RedisContext, rKey, val).Val() {
			val = x.RandomInt(min, max)
		} else {
			redis.Redis.SAdd(redis.RedisContext, rKey, val)
			break
		}
	}
	return val
}

// RemoveFixed 移出指定序列号
func (*serial) RemoveFixed(tag string, val string) {
	rKey := redis.RedisDataPrefix + "Serial:Rdm:Fixed:" + tag
	redis.Redis.SRem(redis.RedisContext, rKey, val)
}
