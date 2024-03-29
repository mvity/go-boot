/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package core

import (
	"github.com/mvity/go-boot/internal/api"
	"github.com/mvity/go-boot/internal/conf"
	"github.com/mvity/go-boot/internal/dao"
	"github.com/mvity/go-boot/internal/dao/mysql"
	"github.com/mvity/go-boot/internal/dao/redis"
	"github.com/mvity/go-boot/internal/job"
	"github.com/mvity/go-boot/internal/logs"
	"github.com/mvity/go-boot/internal/ws"
	"log"
	"os"
)

// Boot 启动入口
func Boot(_api, _job, _wss bool, configFilePath string, port int) {
	log.Println("SYS [INFO] GoBoot version: ", Version)
	log.Println("SYS [INFO] GoBoot Run Module: ", "API [", _api, "]", "Job [", _job, "]", "WebSocket [", _wss, "]")
	if err := conf.InitConfig(configFilePath); err != nil {
		log.Panicf("Init Config error, cause: %v\n", err)
	}
	if err := logs.InitLogger(); err != nil {
		log.Panicf("Init Logger error, cause: %v\n", err)
	}
	if err := dbs.InitMySQL(); err != nil {
		log.Panicf("Init MySQL error, cause: %v\n", err)
	}
	if err := rds.InitRedis(); err != nil {
		log.Panicf("Init Redis error, cause: %v\n", err)
	}

	if port > 0 {
		if _api {
			conf.Config.Port.ApiPort = port
		}
		if _wss {
			conf.Config.Port.WebSocketPort = port
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
		if err := ws.InitWssService(); err != nil {
			log.Panicf("Init WSS Service error, cause: %v\n", err)
		} else {
			log.Println("Start WSS Service success.")
		}
	}

}

// InitProject 初始化项目
func InitProject(configFilePath string) {
	log.Println("WebsocketServer version: ", Version)
	log.Println("Now init project data.")
	if err := conf.InitConfig(configFilePath); err != nil {
		log.Panicf("Init Config error, cause: %v\n", err)
	}
	if err := os.MkdirAll(conf.Config.App.LogPath, os.ModePerm); err != nil {
		log.Panicf("Init LogPath error, cause: %v\n", err)
	}
	if err := dao.InitMySQLDatabase(); err != nil {
		log.Panicf("Init InitMySQLDatabase error, cause: %v\n", err)
	}
	if err := dbs.InitMySQL(); err != nil {
		log.Panicf("Init MySQL error, cause: %v\n", err)
	}
	if err := dao.InitMySQLTable(); err != nil {
		log.Panicf("Init InitMySQLTable error, cause: %v\n", err)
	}
	log.Println("Init project datas success.")
}
