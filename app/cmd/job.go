package cmd

import (
	"2049links-server/job"
	"github.com/sunmi-OS/gocore/v2/utils/closes"
	"github.com/urfave/cli/v2"
)

// Job cmd 任务相关
var Job = &cli.Command{
	Name:    "job",
	Aliases: []string{"j"},
	Usage:   "job",
	Subcommands: []*cli.Command{
		{
			Name:   "InitUser",
			Usage:  "初始化默认用户",
			Action: InitUser,
		},
	},
}

func InitUser(c *cli.Context) error {
	// 初始化配置
	initConf()

	defer closes.Close()
	initDB()
	job.InitUser()
	return nil
}
