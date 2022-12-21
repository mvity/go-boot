package rds

import (
	"fmt"
	"github.com/mvity/go-boot/internal/app"
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

	ukey := RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	tkey := RedisDataPrefix + "Token:T:"

	// 清空当前已登录Token信息
	if otv := t.GetToken(userId, tag); otv != "" {
		if cmd := Redis.Set(RedisContext, tkey+otv, "Killed", 10*time.Minute); cmd.Err() != nil || cmd.Val() != "OK" {
			panic(&app.RedisError{Message: "清空当前登录Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
		}
	}
	// 设置Token对应的UserId
	if cmd := Redis.Set(RedisContext, tkey+tval, strconv.FormatUint(userId, 10), time.Duration(minute)*time.Minute); cmd.Err() != nil || cmd.Val() != "OK" {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	}
	// 设置UserID和APP 对应的Token
	if cmd := Redis.HSet(RedisContext, ukey, tag, tval); cmd.Err() != nil {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	}
	// 读取Ukey过期剩余秒数
	if cmd := Redis.TTL(RedisContext, ukey); cmd.Err() != nil {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	} else {
		ttl := cmd.Val()
		if ttl.Seconds() < float64(minute*60) {
			if scmd := Redis.Expire(RedisContext, ukey, time.Duration(minute)*time.Minute); scmd.Err() != nil {
				panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", scmd.Args())})
			}
		}
	}

	return tval
}

// LogoutWithUserId 注销指定用户
func (t *token) LogoutWithUserId(userId uint64) {
	ukey := RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	for _, tval := range Redis.HGetAll(RedisContext, ukey).Val() {
		t.clearToken(tval)
	}
	Redis.Del(RedisContext, ukey)
}

// LogoutWithToken 注销指定Token
func (t *token) LogoutWithToken(token string) {
	t.clearToken(token)
}

// GetToken 获取指定用户Token值
func (t *token) GetToken(userId uint64, tag string) string {
	ukey := RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	return Redis.HGet(RedisContext, ukey, tag).Val()
}

// GetUserId 获取指定Token用户ID
func (t *token) GetUserId(token string) int64 {
	if strings.EqualFold(token, "guest") {
		return 0
	}
	tkey := RedisDataPrefix + "Token:T:" + token
	idStr := Redis.Get(RedisContext, tkey).Val()
	if idStr == "" {
		return 0
	}
	if idStr == "Killed" {
		return -1
	}
	if val, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		if minute, err := strconv.ParseInt(strings.Split(token, "_")[1], 10, 64); err == nil && minute > 0 {
			Redis.Expire(RedisContext, tkey, time.Duration(minute)*time.Minute)
		}
		return int64(val)
	}
	return 0
}

func (t *token) clearToken(token string) {
	userId := t.GetUserId(token)
	if userId > 0 {
		ukey := RedisDataPrefix + "Token:U:" + strconv.FormatUint(uint64(userId), 10)
		for tkey, tval := range Redis.HGetAll(RedisContext, ukey).Val() {
			if tval == token {
				Redis.HDel(RedisContext, ukey, tkey)
			}
		}
	}
	tkey := RedisDataPrefix + "Token:T:" + token
	Redis.Del(RedisContext, tkey)
}
