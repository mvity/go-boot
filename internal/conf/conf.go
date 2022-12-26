/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package conf

type appConf struct {
	Debug   bool   `yaml:"debug"`
	ApiRoot string `yaml:"api"`
	LogPath string `yaml:"log"`
}

type portConf struct {
	ApiPort       int `yaml:"api"`
	WebSocketPort int `yaml:"ws"`
}

type dataConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	MaxConn  int    `yaml:"max-conn"`
	MaxIdle  int    `yaml:"max-idle"`
	Timeout  int    `yaml:"timeout"`
}

type config struct {
	App  appConf  `yaml:"app"`
	Port portConf `yaml:"port"`
	Data struct {
		MySQL dataConf `yaml:"mysql"`
		Redis dataConf `yaml:"redis"`
	} `yaml:"data"`
}

var Config *config
