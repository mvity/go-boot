/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
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

type mysqlConf struct {
	DSN             string `yaml:"dsn"`
	Database        string `yaml:"database"`
	MaxOpen         int    `yaml:"max-open"`
	MaxIdle         int    `yaml:"max-idle"`
	MaxIdleTime     int    `yaml:"max-idle-time"`
	MaxConnLifetime int    `yaml:"max-conn-lifetime"`
}

type redisConf struct {
	Addr     string `yaml:"addr"`
	Database int    `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	MinIdle  int    `yaml:"min-idle"`
	MaxIdle  int    `yaml:"max-idle"`
	Prefix   string `yaml:"prefix"`
}

type config struct {
	App  appConf  `yaml:"app"`
	Port portConf `yaml:"port"`
	Data struct {
		MySQL mysqlConf `yaml:"mysql"`
		Redis redisConf `yaml:"redis"`
	} `yaml:"data"`
}

var Config *config
