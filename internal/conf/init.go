/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package conf

import (
	"github.com/mvity/go-box/x"
	"gopkg.in/yaml.v2"
	"os"
)

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
