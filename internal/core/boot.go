package core

import (
	"github.com/mvity/go-boot/internal/api"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-boot/internal/dao/mysql"
	"github.com/mvity/go-boot/internal/dao/redis"
	"github.com/mvity/go-boot/internal/job"
	"github.com/mvity/go-boot/internal/wss"
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
	if err := mysql.InitMySQL(); err != nil {
		log.Panicf("Init MySQL error, cause: %v\n", err)
	}
	if err := redis.InitRedis(); err != nil {
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
		} else {
			log.Println("Start API Service success.")
		}
	}

	if _job {
		if err := job.InitJobService(); err != nil {
			log.Panicf("Init Job Service error, cause: %v\n", err)
		} else {
			log.Println("Start Job Service success.")
		}
	}

	if _wss {
		if err := wss.InitWssService(); err != nil {
			log.Panicf("Init WSS Service error, cause: %v\n", err)
		} else {
			log.Println("Start WSS Service success.")
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
	if err := mysql.InitMySQL(); err != nil {
		log.Panicf("Init MySQL error, cause: %v\n", err)
	}
	if err := redis.InitRedis(); err != nil {
		log.Panicf("Init Redis error, cause: %v\n", err)
	}
	log.Println("Init project datas success.")
}
