package core

import (
	"github.com/mvity/go-quickstart/internal/api"
	"github.com/mvity/go-quickstart/internal/app"
	"github.com/mvity/go-quickstart/internal/dao"
	"log"
)

// Boot 启动入口
func Boot(_api, _job, _wss bool, config string, port int) {
	log.Println("Server version: ", Version)
	log.Println("Server Run Module: ", "API [", _api, "]", "Job [", _job, "]", "WebSocket [", _wss, "]")
	if err := app.InitConfig(config); err != nil {
		log.Panicf("Init Config error, cause: %v\n", err)
	}
	if err := app.InitLogger(); err != nil {
		log.Panicf("Init Logger error, cause: %v\n", err)
	}
	if err := dao.InitMySQL(); err != nil {
		log.Panicf("Init MySQL error, cause: %v\n", err)
	}
	if err := dao.InitRedis(); err != nil {
		log.Panicf("Init Redis error, cause: %v\n", err)
	}

	if port > 0 {
		if _api {
			app.Config.Port.ApiPort = port
		}
		if _wss {
			app.Config.Port.WebSocketPort = port
		}
	}
	if _api {
		if err := api.InitApiService(); err != nil {
			log.Panicf("Init API Service error, cause: %v\n", err)
		}
	}

}

// InitProject 初始化项目
func InitProject(config string) {
	log.Println("Server version: ", Version)
	log.Println("Now init project datas.")
	if err := app.InitConfig(config); err != nil {
		log.Panicf("Init Config error, cause: %v\n", err)
	}
	if err := dao.InitMySQL(); err != nil {
		log.Panicf("Init MySQL error, cause: %v\n", err)
	}
	if err := dao.InitRedis(); err != nil {
		log.Panicf("Init Redis error, cause: %v\n", err)
	}
	log.Println("Init project datas success.")
}
