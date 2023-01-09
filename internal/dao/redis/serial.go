/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package rds

import (
	"errors"
	"fmt"
	"github.com/mvity/go-box/x"
	"strconv"
	"strings"
	"time"
)

type serial struct{}

var Serial serial

// Next 获取下一个序列号
func (*serial) Next(tag string, prifix string, size int, expireMinutes int64) string {
	rKey := RedisDataPrefix + "Serial:Seq:" + tag + ":" + prifix
	max, _ := strconv.ParseInt(strings.Repeat("9", size), 10, 64)
	init := Redis.Exists(RedisContext, rKey).Val() > 0
	val := Redis.IncrBy(RedisContext, rKey, 1).Val()
	if val > max {
		panic(errors.New("out of maximum serial number"))
	}
	if !init {
		Redis.Expire(RedisContext, rKey, time.Duration(expireMinutes)*time.Minute)
	}
	return fmt.Sprintf("%0"+strconv.FormatInt(int64(size), 10)+"d", val)
}

// Random 获取下一个随机序列号
func (*serial) Random(tag string, prifix string, size int, expireMinutes int64) string {
	rKey := RedisDataPrefix + "Serial:Rdm:" + tag + ":" + prifix
	init := Redis.Exists(RedisContext, rKey).Val() > 0
	val := x.RandomString(size, false, true)
	for i := 0; i < size; i++ {
		if Redis.SIsMember(RedisContext, rKey, val).Val() {
			val = x.RandomString(size, false, true)
		} else {
			Redis.SAdd(RedisContext, rKey, val)
			break
		}
	}
	if !init {
		Redis.Expire(RedisContext, rKey, time.Duration(expireMinutes)*time.Minute)
	}
	return val
}

// RandomFixed 获取下一个随机序列号
func (*serial) RandomFixed(tag string, min int64, max int64) int64 {
	rKey := RedisDataPrefix + "Serial:Rdm:Fixed:" + tag
	val := x.RandomInt(min, max)
	for i := min; i < max; i++ {
		if Redis.SIsMember(RedisContext, rKey, val).Val() {
			val = x.RandomInt(min, max)
		} else {
			Redis.SAdd(RedisContext, rKey, val)
			break
		}
	}
	return val
}

// RemoveFixed 移出指定序列号
func (*serial) RemoveFixed(tag string, val string) {
	rKey := RedisDataPrefix + "Serial:Rdm:Fixed:" + tag
	Redis.SRem(RedisContext, rKey, val)
}
