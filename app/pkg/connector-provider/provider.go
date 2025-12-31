package connector_provider

import (
	"os"

	"github.com/markbates/goth/providers/github"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"

	"forgeturl-server/pkg/core"

	"github.com/forgeturl/redistore"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// getConfig 获取配置值，优先从 viper 配置读取，若为空则从环境变量读取
func getConfig(viperKey, envKey string) string {
	if val := viper.C.GetString(viperKey); val != "" {
		return val
	}
	return os.Getenv(envKey)
}

func Init() {
	maxAge := 86400 * 30 // 30 days
	isProd := false
	store, _ := redistore.NewRediStore(10, "tcp", ":6379", "", "", []byte("fg-key"))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.SetMaxAge(maxAge)

	gothic.Store = store
	providers := []goth.Provider{}

	// Google provider
	googleKey := getConfig("keys.GOOGLE_KEY", "GOOGLE_KEY")
	googleSecret := getConfig("keys.GOOGLE_SECRET", "GOOGLE_SECRET")
	if googleKey != "" && googleSecret != "" {
		providers = append(providers, google.New(googleKey, googleSecret, core.FillDomain("/auth/callback/google")))
		glog.InfoF("google provider inited")
	}

	// GitHub provider
	githubKey := getConfig("keys.GITHUB_KEY", "GITHUB_KEY")
	githubSecret := getConfig("keys.GITHUB_SECRET", "GITHUB_SECRET")
	if githubKey != "" && githubSecret != "" {
		providers = append(providers, github.New(githubKey, githubSecret, core.FillDomain("/auth/callback/github")))
		glog.InfoF("github provider inited")
	}

	if len(providers) == 0 {
		panic("no oauth provider configured")
	}

	goth.UseProviders(providers...)
}
