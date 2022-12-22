package job

// InitJobService 启动JobTask服务
func InitJobService() error {
	go Executor.Start()
	go Fiexd.Start()
	select {}
}
