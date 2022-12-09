package core

import (
	"github.com/mvity/go-quickstart/internal/app"
	"log"
)

// Boot 启动入口
func Boot(api, job, wss bool, config string, port int8) {
	log.Println("Server version: ", Version)
	log.Println("Server Run Module: ", "API [", api, "]", "Job [", job, "]", "WebSocket [", wss, "]")
	if err := app.InitConfig(config); err != nil {
		log.Panicf("Init Config error, cause: %v\n", err)
	}
	if err := app.InitLogger(); err != nil {
		log.Panicf("Init Logger error, cause: %v\n", err)
	}

}

// InitProject 初始化项目
func InitProject(config string) {
	log.Println("Server version: ", Version)
	log.Println("Now init project datas.")

	if err := app.InitConfig(config); err != nil {
		log.Panicf("Init Config error, cause: %v\n", err)
	}
	log.Println("Init project datas success.")
}
