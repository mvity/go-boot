package job

import (
	"context"
	"fmt"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao/dbs"
	rds2 "github.com/mvity/go-boot/internal/dao/rds"
	"github.com/mvity/go-box/x"
	"gorm.io/gorm"

	"time"
)

type executor struct {
	ticker *time.Ticker
}

// Executor 动态任务执行器
var Executor executor

// Start 启动执行器
func (e *executor) Start() {
	e.ticker = time.NewTicker(time.Second * time.Duration(1))
	for {
		go e.handle(time.Now())
		<-e.ticker.C
	}
}

// Stop 停止执行器
func (e *executor) Stop() {
	e.ticker.Stop()
}

// 处理任务
func (e *executor) handle(now time.Time) {
	tasks := rds2.Task.GetTasks(now)
	for _, task := range tasks {
		app.LogSysInfo(fmt.Sprintf("Exec task %v", task), nil)
		node, err := x.JsonFromStringE(task)
		if err != nil {
			app.LogSysInfo("Parse Task json error ", err)
			continue
		}
		db := dbs.MySQL.WithContext(context.Background())
		switch node.Name("tag").String() {
		case "TestTask":
			go e.doExecTestTask(db, now, node.Name("id").String(), node.Name("info").String())
		default:
			rds2.Task.AddTask(time.Now().Add(time.Minute*5), node.Name("tag").String(), node.Name("info").String())
		}
	}
}

// 执行任务
func (*executor) doExecTestTask(db *gorm.DB, now time.Time, id string, info string) {
	defer func() {
		if err := recover(); err != nil {
			app.LogSysInfo(fmt.Sprintf("Exec doExecTestTask Error on %v , info: %v \n", x.FormatDateTime(now), info), err.(error))
		}
	}()
	if rds2.Task.IsHandled(now, id) {
		return
	}
	if ok := rds2.Locker.Lock("Task:Handle:"+id, 300); !ok {
		return
	} else {
		defer rds2.Locker.UnLock("Task:Handle:" + id)
	}
	rds2.Task.AddHandled(now, id)

	fmt.Println("Exec TestTask biz....")

}
