/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package job

import (
	"context"
	"fmt"
	"github.com/mvity/go-boot/internal/dao/mysql"
	rds "github.com/mvity/go-boot/internal/dao/redis"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-box/x"
	"gorm.io/gorm"

	"time"
)

type task struct {
	ticker *time.Ticker
}

// Task 动态任务执行器
var Task task

// Start 启动执行器
func (e *task) Start() {
	e.ticker = time.NewTicker(time.Second * time.Duration(1))
	for {
		go e.handle(time.Now())
		<-e.ticker.C
	}
}

// Stop 停止执行器
func (e *task) Stop() {
	e.ticker.Stop()
}

// 处理任务
func (e *task) handle(now time.Time) {
	tasks := rds.Task.GetTasks(now)
	for _, task := range tasks {
		logs.LogJobInfo(fmt.Sprintf("Exec task %v", task), nil)
		node, err := x.JsonFromStringE(task)
		if err != nil {
			logs.LogJobInfo("Parse Task json error ", err)
			continue
		}
		db := dbs.MySQL.WithContext(context.Background())
		switch node.Name("tag").String() {
		case "TestTask":
			go e.doExecTestTask(db, now, node.Name("id").String(), node.Name("info").String())
		default:
			rds.Task.AddTask(time.Now().Add(time.Minute*5), node.Name("tag").String(), node.Name("info").String())
		}
	}
}

// 执行任务
func (*task) doExecTestTask(db *gorm.DB, now time.Time, id string, info string) {
	defer func() {
		if err := recover(); err != nil {
			logs.LogJobInfo(fmt.Sprintf("Exec doExecTestTask Error on %v , info: %v \n", x.FormatDateTime(now), info), err.(error))
		}
	}()
	if rds.Task.IsHandled(now, id) {
		return
	}
	if ok := rds.Locker.Lock("Task:Handle:"+id, 300); !ok {
		return
	} else {
		defer rds.Locker.UnLock("Task:Handle:" + id)
	}
	rds.Task.AddHandled(now, id)

	fmt.Println("Exec TestTask biz....")

}
