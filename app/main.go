package main

import (
	"os"

	"forgeturl-server/cmd"
	"forgeturl-server/conf"

	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	// 打印Banner
	utils.PrintBanner(conf.ProjectName)
	// 配置cli参数
	app := cli.NewApp()
	app.Name = conf.ProjectName
	app.Usage = conf.ProjectName
	app.Version = conf.ProjectVersion

	// 指定命令运行的函数
	app.Commands = []*cli.Command{
		cmd.Api,
		cmd.Job,
	}

	// 启动cli
	if err := app.Run(os.Args); err != nil {
		glog.ErrorF("Failed to start application: %v", err)
	}
}
