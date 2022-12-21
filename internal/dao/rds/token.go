package rds

import (
	"fmt"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao"
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

	ukey := dao.RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	tkey := dao.RedisDataPrefix + "Token:T:"

	// 清空当前已登录Token信息
	if otv := t.GetToken(userId, tag); otv != "" {
		if cmd := dao.Redis.Set(dao.MySQLContext, tkey+otv, "Killed", 10*time.Minute); cmd.Err() != nil || cmd.Val() != "OK" {
			panic(&app.RedisError{Message: "清空当前登录Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
		}
	}
	// 设置Token对应的UserId
	if cmd := dao.Redis.Set(dao.MySQLContext, tkey+tval, strconv.FormatUint(userId, 10), time.Duration(minute)*time.Minute); cmd.Err() != nil || cmd.Val() != "OK" {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	}
	// 设置UserID和APP 对应的Token
	if cmd := dao.Redis.HSet(dao.MySQLContext, ukey, tag, tval); cmd.Err() != nil {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	}
	// 读取Ukey过期剩余秒数
	if cmd := dao.Redis.TTL(dao.MySQLContext, ukey); cmd.Err() != nil {
		panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", cmd.Args())})
	} else {
		ttl := cmd.Val()
		if ttl.Seconds() < float64(minute*60) {
			if scmd := dao.Redis.Expire(dao.MySQLContext, ukey, time.Duration(minute)*time.Minute); scmd.Err() != nil {
				panic(&app.RedisError{Message: "设置Token失败", Command: fmt.Sprintf("%v", scmd.Args())})
			}
		}
	}

	return tval
}

// LogoutWithUserId 注销指定用户
func (t *token) LogoutWithUserId(userId uint64) {
	ukey := dao.RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	for _, tval := range dao.Redis.HGetAll(dao.MySQLContext, ukey).Val() {
		t.clearToken(tval)
	}
	dao.Redis.Del(dao.MySQLContext, ukey)
}

// LogoutWithToken 注销指定Token
func (t *token) LogoutWithToken(token string) {
	t.clearToken(token)
}

// GetToken 获取指定用户Token值
func (t *token) GetToken(userId uint64, tag string) string {
	ukey := dao.RedisDataPrefix + "Token:U:" + strconv.FormatUint(userId, 10)
	return dao.Redis.HGet(dao.MySQLContext, ukey, tag).Val()
}

// GetUserId 获取指定Token用户ID
func (t *token) GetUserId(token string) int64 {
	if strings.EqualFold(token, "guest") {
		return 0
	}
	tkey := dao.RedisDataPrefix + "Token:T:" + token
	idStr := dao.Redis.Get(dao.MySQLContext, tkey).Val()
	if idStr == "" {
		return 0
	}
	if idStr == "Killed" {
		return -1
	}
	if val, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		if minute, err := strconv.ParseInt(strings.Split(token, "_")[1], 10, 64); err == nil && minute > 0 {
			dao.Redis.Expire(dao.MySQLContext, tkey, time.Duration(minute)*time.Minute)
		}
		return int64(val)
	}
	return 0
}

func (t *token) clearToken(token string) {
	userId := t.GetUserId(token)
	if userId > 0 {
		ukey := dao.RedisDataPrefix + "Token:U:" + strconv.FormatUint(uint64(userId), 10)
		for tkey, tval := range dao.Redis.HGetAll(dao.MySQLContext, ukey).Val() {
			if tval == token {
				dao.Redis.HDel(dao.MySQLContext, ukey, tkey)
			}
		}
	}
	tkey := dao.RedisDataPrefix + "Token:T:" + token
	dao.Redis.Del(dao.MySQLContext, tkey)
}
