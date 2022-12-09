package app

import (
	"github.com/mvity/go-box/x"
	"gopkg.in/yaml.v2"
	"os"
)

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

// InitConfig 初始化配置文件
func InitConfig(configFilePath string) error {
	pwd, _ := os.Getwd()
	cfp := x.StringDefaultIfBlank(configFilePath, pwd+"./configs/conf.yaml")
	bytes, err := os.ReadFile(cfp)
	if err != nil {
		return err
	}
	Config = new(config)
	return yaml.Unmarshal(bytes, Config)
}
