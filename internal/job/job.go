package job

import "github.com/mvity/go-boot/internal/app"

// InitJobService 启动JobTask服务
func InitJobService(daemon bool) error {
	go Executor.Start()
	go Fiexd.Start()
	app.LogSysInfo("Start JOB service success", nil)
	if daemon {
		select {}
	}
	return nil
}
