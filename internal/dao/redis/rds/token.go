package rds

import (
	"fmt"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao/redis"
	"github.com/mvity/go-box/x"
	"strconv"
	"strings"
	"time"
)

type token struct{}

var Token token

// Login 生成新的登录信息
func (t *token) Login(userId uint64, tag string, minute int) string {

	md5v := x.MD5String(strconv.FormatUint(userId, 10) + strconv.FormatInt(time.Now().Unix(), 10) + x.RandomString(32, true, true))

	suffix := strconv.FormatInt(int64(minute), 10)

	tval := strings.ToUpper(md5v+x.RandomString(16, true, true)) + "_" + suffix

	ukey := redis.RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	tkey := redis.RedisDataPrefix + "Token:T:"

	// 清空当前已登录Token信息
	if otv := t.GetToken(userId, tag); otv != "" {
		if cmd := redis.Redis.Set(redis.RedisContext, tkey+otv, "Killed", 10*time.Minute); cmd.Err() != nil || cmd.Val() != "OK" {
			panic(&app.RedisError{Message: "清空当前登录Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
		}
	}
	// 设置Token对应的UserId
	if cmd := redis.Redis.Set(redis.RedisContext, tkey+tval, strconv.FormatUint(userId, 10), time.Duration(minute)*time.Minute); cmd.Err() != nil || cmd.Val() != "OK" {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	}
	// 设置UserID和APP 对应的Token
	if cmd := redis.Redis.HSet(redis.RedisContext, ukey, tag, tval); cmd.Err() != nil {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	}
	// 读取Ukey过期剩余秒数
	if cmd := redis.Redis.TTL(redis.RedisContext, ukey); cmd.Err() != nil {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	} else {
		ttl := cmd.Val()
		if ttl.Seconds() < float64(minute*60) {
			if scmd := redis.Redis.Expire(redis.RedisContext, ukey, time.Duration(minute)*time.Minute); scmd.Err() != nil {
				panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", scmd.Args())})
			}
		}
	}

	return tval
}

// LogoutWithUserId 注销指定用户
func (t *token) LogoutWithUserId(userId uint64) {
	ukey := redis.RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	for _, tval := range redis.Redis.HGetAll(redis.RedisContext, ukey).Val() {
		t.clearToken(tval)
	}
	redis.Redis.Del(redis.RedisContext, ukey)
}

// LogoutWithToken 注销指定Token
func (t *token) LogoutWithToken(token string) {
	t.clearToken(token)
}

// GetToken 获取指定用户Token值
func (t *token) GetToken(userId uint64, tag string) string {
	ukey := redis.RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	return redis.Redis.HGet(redis.RedisContext, ukey, tag).Val()
}

// GetUserId 获取指定Token用户ID
func (t *token) GetUserId(token string) int64 {
	if strings.EqualFold(token, "guest") {
		return 0
	}
	tkey := redis.RedisDataPrefix + "Token:T:" + token
	idStr := redis.Redis.Get(redis.RedisContext, tkey).Val()
	if idStr == "" {
		return 0
	}
	if idStr == "Killed" {
		return -1
	}
	if val, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		if minute, err := strconv.ParseInt(strings.Split(token, "_")[1], 10, 64); err == nil && minute > 0 {
			redis.Redis.Expire(redis.RedisContext, tkey, time.Duration(minute)*time.Minute)
		}
		return int64(val)
	}
	return 0
}

func (t *token) clearToken(token string) {
	userId := t.GetUserId(token)
	if userId > 0 {
		ukey := redis.RedisDataPrefix + "Token:U:" + strconv.FormatUint(uint64(userId), 10)
		for tkey, tval := range redis.Redis.HGetAll(redis.RedisContext, ukey).Val() {
			if tval == token {
				redis.Redis.HDel(redis.RedisContext, ukey, tkey)
			}
		}
	}
	tkey := redis.RedisDataPrefix + "Token:T:" + token
	redis.Redis.Del(redis.RedisContext, tkey)
}
