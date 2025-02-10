package core

import (
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"os"
)

var zone string

const (
	TestEnv = "test"
	ProdEnv = "prod"
)

// IsTestEnv 是否是测试环境
// 线上环境必须赋值ZONE环境变量
func IsTestEnv() bool {
	if zone == "" {
		zone = os.Getenv("ZONE")
		if zone == "" {
			zone = TestEnv
		}
	}

	if zone == TestEnv {
		return true
	}
	return false
}

func IsProdEnv() bool {
	return !IsTestEnv()
}

func FillDomain(path string) string {
	if IsTestEnv() {
		return "http://" + viper.C.GetString("base.domain") + path
	}
	return "https://" + viper.C.GetString("base.domain") + path
}
