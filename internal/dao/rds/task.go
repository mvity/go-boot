package rds

import (
	"errors"
	"fmt"
	v9 "github.com/go-redis/redis/v9"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao"
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

	rkey := dao.RedisDataPrefix + "Tasks:Run:" + exec.Format("20060102")
	if err := dao.Redis.ZAdd(dao.MySQLContext, rkey, v9.Z{
		Score:  float64(exec.Unix()),
		Member: x.JsonToString(args),
	}).Err(); err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v add error", info, exec.Format(time.RFC3339))))
	}
	dao.Redis.ExpireAt(dao.MySQLContext, rkey, exec.Add(time.Duration(48)*time.Hour))
}

// GetTasks 获取指定时间的任务数据
func (*task) GetTasks(now time.Time) []string {
	rkey := dao.RedisDataPrefix + "Tasks:Run:" + now.Format("20060102")

	score := strconv.FormatInt(now.Unix(), 10)

	infos := dao.Redis.ZRangeByScore(dao.MySQLContext, rkey, &v9.ZRangeBy{
		Min: "0",
		Max: score,
	}).Val()
	dao.Redis.ZRemRangeByScore(dao.MySQLContext, rkey, "0", score)

	return infos
}

// AddHandled 添加已执行任务信息
func (*task) AddHandled(now time.Time, id string) {
	rkey := dao.RedisDataPrefix + "Tasks:End:" + now.Format("20060102")
	if err := dao.Redis.SAdd(dao.MySQLContext, rkey, id).Err(); err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v add handled error", id, now.Format(time.RFC3339))))
	}
	dao.Redis.ExpireAt(dao.MySQLContext, rkey, now.Add(time.Duration(48)*time.Hour))
}

// IsHandled 指定任务是否已执行
func (*task) IsHandled(now time.Time, id string) bool {
	rkey := dao.RedisDataPrefix + "Tasks:End:" + now.Format("20060102")
	exists, err := dao.Redis.SIsMember(dao.MySQLContext, rkey, id).Result()
	if err != nil {
		panic(errors.New(fmt.Sprintf("Task [%v] @ %v check handled error", id, now.Format(time.RFC3339))))
	}
	return exists
}
