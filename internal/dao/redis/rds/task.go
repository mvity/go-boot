package rds

import (
	"errors"
	"fmt"
	v9 "github.com/go-redis/redis/v9"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao/redis"
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

	rkey := redis.DataPrefix + "Tasks:Run:" + exec.Format("20060102")
	if err := redis.Redis.ZAdd(redis.Context, rkey, v9.Z{
		Score:  float64(exec.Unix()),
		Member: x.JsonToString(args),
	}).Err(); err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v add error", info, exec.Format(time.RFC3339))))
	}
	redis.Redis.ExpireAt(redis.Context, rkey, exec.Add(time.Duration(48)*time.Hour))
}

// GetTasks 获取指定时间的任务数据
func (*task) GetTasks(now time.Time) []string {
	rkey := redis.DataPrefix + "Tasks:Run:" + now.Format("20060102")

	score := strconv.FormatInt(now.Unix(), 10)

	infos := redis.Redis.ZRangeByScore(redis.Context, rkey, &v9.ZRangeBy{
		Min: "0",
		Max: score,
	}).Val()
	redis.Redis.ZRemRangeByScore(redis.Context, rkey, "0", score)

	return infos
}

// AddHandled 添加已执行任务信息
func (*task) AddHandled(now time.Time, id string) {
	rkey := redis.DataPrefix + "Tasks:End:" + now.Format("20060102")
	if err := redis.Redis.SAdd(redis.Context, rkey, id).Err(); err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v add handled error", id, now.Format(time.RFC3339))))
	}
	redis.Redis.ExpireAt(redis.Context, rkey, now.Add(time.Duration(48)*time.Hour))
}

// IsHandled 指定任务是否已执行
func (*task) IsHandled(now time.Time, id string) bool {
	rkey := redis.DataPrefix + "Tasks:End:" + now.Format("20060102")
	exists, err := redis.Redis.SIsMember(redis.Context, rkey, id).Result()
	if err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v check handled error", id, now.Format(time.RFC3339))))
	}
	return exists
}
