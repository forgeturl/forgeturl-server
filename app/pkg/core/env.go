package core

import (
	"fmt"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/utils"
)

func FillDomain(path string) string {
	if utils.IsLocal() {
		return fmt.Sprintf("%s:%s%s", viper.C.GetString("base.domain"), viper.C.GetString("network.ApiServicePort"), path)
	}
	return viper.C.GetString("base.domain") + path
}
