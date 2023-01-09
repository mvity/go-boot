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
	v9 "github.com/go-redis/redis/v9"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-box/x"
	"strconv"
	"time"
)

type task struct{}

// Task 任务管理器
var Task task

// AddTask 添加任务
func (*task) AddTask(exec time.Time, tag string, info string) {
	args := make(map[string]any)
	args["id"] = x.ToString(app.IDWorker.GetID())
	args["tag"] = tag
	args["info"] = info

	rkey := RedisDataPrefix + "Tasks:Run:" + exec.Format("20060102")
	if err := Redis.ZAdd(RedisContext, rkey, v9.Z{
		Score:  float64(exec.Unix()),
		Member: x.JsonToString(args),
	}).Err(); err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v add error", info, exec.Format(time.RFC3339))))
	}
	Redis.ExpireAt(RedisContext, rkey, exec.Add(time.Duration(48)*time.Hour))
}

// GetTasks 获取指定时间的任务数据
func (*task) GetTasks(now time.Time) []string {
	rkey := RedisDataPrefix + "Tasks:Run:" + now.Format("20060102")

	score := strconv.FormatInt(now.Unix(), 10)

	infos := Redis.ZRangeByScore(RedisContext, rkey, &v9.ZRangeBy{
		Min: "0",
		Max: score,
	}).Val()
	Redis.ZRemRangeByScore(RedisContext, rkey, "0", score)

	return infos
}

// AddHandled 添加已执行任务信息
func (*task) AddHandled(now time.Time, id string) {
	rkey := RedisDataPrefix + "Tasks:End:" + now.Format("20060102")
	if err := Redis.SAdd(RedisContext, rkey, id).Err(); err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v add handled error", id, now.Format(time.RFC3339))))
	}
	Redis.ExpireAt(RedisContext, rkey, now.Add(time.Duration(48)*time.Hour))
}

// IsHandled 指定任务是否已执行
func (*task) IsHandled(now time.Time, id string) bool {
	rkey := RedisDataPrefix + "Tasks:End:" + now.Format("20060102")
	exists, err := Redis.SIsMember(RedisContext, rkey, id).Result()
	if err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v check handled error", id, now.Format(time.RFC3339))))
	}
	return exists
}
