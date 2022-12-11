package rds

import (
	"github.com/mvity/go-box/x"
	"github.com/mvity/go-quickstart/internal/dao/redis"
	"time"
)

type smsCaptcha struct{}

// SmsCaptcha 短信验证码
var SmsCaptcha smsCaptcha

// getCaptchInfo 获取验证码信息
func (s *smsCaptcha) getCaptchInfo(mob string, code string) string {
	return x.MD5String(x.MD5String(mob+code) + x.StringReverse(mob))
}

// GenerateCaptch 生成手机号验证码
func (s *smsCaptcha) GenerateCaptch(mob string, minute int) (string, string) {
	// 测试账号
	if mob == "19999999999" {
		return "1234567890", "9999"
	}
	if mob == "18888888888" {
		return "1234567890", "8888"
	}
	if mob == "17777777777" {
		return "1234567890", "7777"
	}
	if mob == "16666666666" {
		return "1234567890", "6666"
	}
	if mob == "15555555555" {
		return "1234567890", "5555"
	}

	code := x.RandomString(4, false, true)
	info := s.getCaptchInfo(mob, code)

	rkey := redis.RedisDataPrefix + "Captch:Sms:" + info
	redis.Redis.IncrBy(redis.RedisContext, rkey, 0)
	redis.Redis.Expire(redis.RedisContext, rkey, time.Duration(minute)*time.Minute)
	return info, code
}

// ValidCaptch 校验验证码
func (s *smsCaptcha) ValidCaptch(mob string, code string, info string) bool {
	// 测试账号
	if mob == "19999999999" && code == "9999" && info == "1234567890" {
		return true
	}
	if mob == "18888888888" && code == "8888" && info == "1234567890" {
		return true
	}
	if mob == "17777777777" && code == "7777" && info == "1234567890" {
		return true
	}
	if mob == "16666666666" && code == "6666" && info == "1234567890" {
		return true
	}
	if mob == "15555555555" && code == "5555" && info == "1234567890" {
		return true
	}

	rkey := redis.RedisDataPrefix + "Captch:Sms:" + info

	if redis.Redis.Exists(redis.RedisContext, rkey).Val() == 0 {
		return false
	}
	if redis.Redis.IncrBy(redis.RedisContext, rkey, 0).Val() >= 6 {
		redis.Redis.Del(redis.RedisContext, rkey)
		return false
	}
	flag := info == s.getCaptchInfo(mob, code)
	if flag {
		redis.Redis.Del(redis.RedisContext, rkey)
	} else {
		redis.Redis.IncrBy(redis.RedisContext, rkey, 1)
	}
	return flag
}
