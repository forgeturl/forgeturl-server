package cmd

import (
	"time"

	"forgeturl-server/pkg/middleware"
	"forgeturl-server/route"

	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore-contrib/smartgzip"
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

	newG := middleware.NewGin()
	//gs := api.NewGinServer(
	//	api.WithServerDebug(!utils.IsRelease()),
	//	api.WithServerHost(viper.C.GetString("network.ApiServiceHost")),
	//	api.WithServerPort(viper.C.GetInt("network.ApiServicePort")),
	//	api.WithServerTimeout(time.Minute*5),
	//	api.WithOpenTrace(false),
	//)
	//gin.DefaultWriter = middleware.NewGlogWriterDebug()
	//gin.DefaultErrorWriter = middleware.NewGlogWriterError()

	addMiddlewares(newG)

	// init route
	route.Routes(newG)
	address := viper.C.GetString("network.ApiServiceHost") + ":" + viper.C.GetString("network.ApiServicePort")
	err := endless.ListenAndServe(address, newG)
	if err != nil {
		return err
	}
	return nil
}

func addMiddlewares(g *gin.Engine) {
	g.Use(middleware.IgnoreNotExistPath())
	g.Use(middleware.MustLocalIp())

	// logging
	g.Use(middleware.ServerLogging(
		middleware.WithSlowThreshold(10000),
		middleware.WithHideReqBodyLogsPath(map[string]bool{
			"/dumplinks/exportBookmarks": true,
		}, true),
	))

	// cors - 自定义配置以支持 withCredentials
	var allowOrigins []string
	switch utils.GetRunTime() {
	case utils.LocalEnv:
		allowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000", "http://127.0.0.1", "http://localhost"}
	case utils.TestEnv:
		allowOrigins = []string{"https://test-api.brightguo.com"}
	case utils.ReleaseEnv:
		allowOrigins = []string{"https://api.brightguo.com"}
	default:
		allowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000", "http://127.0.0.1", "http://localhost"}
	}
	corsConfig := cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Token", "X-Forget-Cookie"},
		ExposeHeaders:    []string{"Content-Length", "X-Token", "X-Forget-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	g.Use(cors.New(corsConfig))

	// gzip
	g.Use(smartgzip.GzipOnly(
		"/space/getPage",
		"/space/getMySpace",
	))
}
