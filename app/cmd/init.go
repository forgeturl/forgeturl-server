package cmd

import (
	"forgeturl-server/conf"
	"forgeturl-server/dal"
	connector_facebook "forgeturl-server/pkg/connector-facebook"
	connector_google "forgeturl-server/pkg/connector-google"
	connector_weixin "forgeturl-server/pkg/connector-weixin"

	"github.com/sunmi-OS/gocore/v2/conf/nacos"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog/zap"
	"github.com/sunmi-OS/gocore/v2/utils"
)

func initConf() {

	switch utils.GetRunTime() {
	case "local":
		nacos.SetLocalConfig(conf.LocalConfig)
	default:
		nacos.SetLocalConfig(conf.OnlConfig)
	}

	vt := nacos.GetViper()
	vt.SetBaseConfig(conf.BaseConfig)
	vt.SetDataIds(conf.ProjectName, "config", "mysql", "redis", "rocketmq")
	// 注册配置更新回调
	vt.SetCallBackFunc(conf.ProjectName, "config", func(namespace, group, dataId, data string) {
		initLog()
	})

	vt.NacosToViper()

}

// initDB 初始化DB服务 （内部方法）
func initDB() {
	dal.Init()
}

// initCache 初始化redis服务 （内部方法）
func initCache() {
	// redis.NewRedis(conf.DemoDb0Redis)
}

// initLog init log
func initLog() {
	zap.SetLogLevel(viper.GetEnvConfig("base.logLevel").String())
}

func initClient() {
	connector_google.Init()
	connector_facebook.Init()
	connector_weixin.Init()
}
