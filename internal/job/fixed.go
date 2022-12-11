package job

import (
	"errors"
	"fmt"
	"github.com/mvity/go-boot/internal/app"
	"github.com/robfig/cron/v3"
)

type fixed struct {
	crond *cron.Cron
}

// Fiexd 定时任务执行器
var Fiexd fixed

// Start 启动执行器
func (f *fixed) Start() {
	f.crond = cron.New()
	f.addFunc()
	f.crond.Start()
}

// Stop 停止执行器
func (f *fixed) Stop() {
	f.crond.Stop()
}

// addFunc 添加任务
func (f *fixed) addFunc() {
	var err error
	var eid cron.EntryID

	if eid, err = f.crond.AddFunc("1 0 * * *", func() {
		fmt.Println("demo")
	}); err != nil {
		panic(errors.New(fmt.Sprintf("[%s] %s %v", "job.Fiexd", "Add [demo] f error", err)))
	} else {
		app.LogSysInfo(fmt.Sprintf("[%s] %s , EntryID: %v", "job.Fiexd", "Add [demo] f success", eid), nil)
	}

}
