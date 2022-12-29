/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"github.com/mvity/go-boot/internal/core"
	"github.com/spf13/cobra"
)

func main() {
	Execute()
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(jobCmd)
	rootCmd.AddCommand(wssCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(verCmd)

	rootCmd.PersistentFlags().StringVarP(&conf, "conf", "c", "", "项目配置文件")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 0, "API/WebSocket 运行端口")
}

// 配置文件
var conf string

// 运行端口
var port int

// 启动主命令
var rootCmd = &cobra.Command{
	Use:   "gboot",
	Short: "一个快速开发Golang应用的模板项目",
	Long:  "Go Quickstart，支持API/Job/WebSocket等方式的项目框架",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Go Quickstart Version " + core.Version + ": no command specified")
		fmt.Println("Try `gboot -h` or `gboot --help` for more information.")
	},
}

// API接口服务启动命令
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "API接口服务",
	Long:  "REST API接口服务",
	Run: func(cmd *cobra.Command, args []string) {
		core.Boot(true, false, false, conf, port)
	},
}

// Job任务服务启动命令
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Job任务服务",
	Long:  "定时Job任务服务",
	Run: func(cmd *cobra.Command, args []string) {
		core.Boot(false, true, false, conf, port)
	},
}

// WebSocket服务启动命令
var wssCmd = &cobra.Command{
	Use:   "ws",
	Short: "WebSocket 服务",
	Long:  "WebSocket 服务",
	Run: func(cmd *cobra.Command, args []string) {
		core.Boot(false, false, true, conf, port)
	},
}

// 初始化项目运行必备条件
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目数据",
	Long:  "同步数据库结构，刷新项目初始化配置信息",
	Run: func(cmd *cobra.Command, args []string) {
		core.InitProject(conf)
	},
}

// 查看当前项目版本
var verCmd = &cobra.Command{
	Use:   "version",
	Short: "查看当前项目版本",
	Long:  "查看当前项目版本",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Server Version : " + core.Version)
	},
}
