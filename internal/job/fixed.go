package job

import (
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
	//var err error
	//var eid cron.EntryID

	//eid, err = f.crond.AddFunc("1 0 * * *", services.SysService.InitPlatSummaryData)
	//if err != nil {
	//	panic(errors.New(fmt.Sprintf("[%s] %s %v", "fiexdJob", "Add [InitPlatSummaryData] f error", err)))
	//} else {
	//	core.AddLog.SysInfo(fmt.Sprintf("[%s] %s , EntryID: %v", "fiexdJob", "Add [InitPlatSummaryData] f success", eid))
	//}
	//
	//eid, err = f.crond.AddFunc("0 4 * * *", services.SysService.SyncDistrict)
	//if err != nil {
	//	panic(errors.New(fmt.Sprintf("[%s] %s %s %v", "fiexdJob", "Start", "Add [SyncDistrict] f error", err)))
	//} else {
	//	core.AddLog.SysInfo(fmt.Sprintf("[%s] %s , EntryID: %v", "fiexdJob", "Add [SyncDistrict] f success", eid))
	//}

	//eid, err = f.crond.AddFunc("0 1 * * *", services.SysService.SyncCarSeries)
	//if err != nil {
	//	panic(errors.New(fmt.Sprintf("[%s] %s %v", "fiexdJob", "Add [SyncCarSeries] f error", err)))
	//} else {
	//	core.AddLog.SysInfo(fmt.Sprintf("[%s] %s , EntryID: %v", "fiexdJob", "Add [SyncCarSeries] f success", eid))
	//}

}
