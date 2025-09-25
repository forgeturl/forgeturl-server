package cmd

import (
	"time"

	"forgeturl-server/pkg/middleware"
	"forgeturl-server/route"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore-contrib/smartgzip"
	"github.com/sunmi-OS/gocore/v2/api"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/urfave/cli/v2"
)

var Api = &cli.Command{
	Name:    "api",
	Aliases: []string{"a"},
	Usage:   "api start",
	Subcommands: []*cli.Command{
		{
			Name:   "start",
			Usage:  "开启运行api服务",
			Action: RunApi,
		},
	},
}

func RunApi(c *cli.Context) error {
	// 初始化配置
	initConf()
	initDB()
	initCache()
	initLog()
	initClient()

	gs := api.NewGinServer(
		api.WithServerDebug(!utils.IsRelease()),
		api.WithServerHost(viper.C.GetString("network.ApiServiceHost")),
		api.WithServerPort(viper.C.GetInt("network.ApiServicePort")),
		api.WithServerTimeout(time.Minute*5),
		api.WithOpenTrace(false),
	)
	gin.DefaultWriter = middleware.NewGlogWriterDebug()
	gin.DefaultErrorWriter = middleware.NewGlogWriterError()

	addMiddlewares(gs.Gin)

	// init route
	route.Routes(gs.Gin)
	gs.Start()
	return nil
}

func addMiddlewares(g *gin.Engine) {
	// logging
	g.Use(middleware.ServerLogging(
		middleware.WithSlowThreshold(10000),
		middleware.WithHideReqBodyLogsPath(map[string]bool{
			"/dumplinks/importBookmarks": true,
			"/dumplinks/exportBookmarks": true,
		}, true),
	))

	// cors
	g.Use(cors.Default())

	// gzip
	g.Use(smartgzip.GzipOnly(
		"/space/getPage",
		"/space/getMySpace",
	))
}
