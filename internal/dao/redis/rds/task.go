package rds

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/mvity/go-boot/internal/app"
	redis2 "github.com/mvity/go-boot/internal/dao/redis"
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

	rkey := redis2.RedisDataPrefix + "Tasks:Run:" + exec.Format("20060102")
	if err := redis2.Redis.ZAdd(redis2.RedisContext, rkey, redis.Z{
		Score:  float64(exec.Unix()),
		Member: x.JsonToString(args),
	}).Err(); err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v add error", info, exec.Format(time.RFC3339))))
	}
	redis2.Redis.ExpireAt(redis2.RedisContext, rkey, exec.Add(time.Duration(48)*time.Hour))
}

// GetTasks 获取指定时间的任务数据
func (*task) GetTasks(now time.Time) []string {
	rkey := redis2.RedisDataPrefix + "Tasks:Run:" + now.Format("20060102")

	score := strconv.FormatInt(now.Unix(), 10)

	infos := redis2.Redis.ZRangeByScore(redis2.RedisContext, rkey, &redis.ZRangeBy{
		Min: "0",
		Max: score,
	}).Val()
	redis2.Redis.ZRemRangeByScore(redis2.RedisContext, rkey, "0", score)

	return infos
}

// AddHandled 添加已执行任务信息
func (*task) AddHandled(now time.Time, id string) {
	rkey := redis2.RedisDataPrefix + "Tasks:End:" + now.Format("20060102")
	if err := redis2.Redis.SAdd(redis2.RedisContext, rkey, id).Err(); err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v add handled error", id, now.Format(time.RFC3339))))
	}
	redis2.Redis.ExpireAt(redis2.RedisContext, rkey, now.Add(time.Duration(48)*time.Hour))
}

// IsHandled 指定任务是否已执行
func (*task) IsHandled(now time.Time, id string) bool {
	rkey := redis2.RedisDataPrefix + "Tasks:End:" + now.Format("20060102")
	exists, err := redis2.Redis.SIsMember(redis2.RedisContext, rkey, id).Result()
	if err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v check handled error", id, now.Format(time.RFC3339))))
	}
	return exists
}
