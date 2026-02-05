package connector_provider

import (
	"fmt"
	"net/http"
	"os"

	"forgeturl-server/pkg/core"

	"github.com/forgeturl/redistore"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/wechat"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
	"golang.org/x/oauth2"
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
	isProd := true
	store, _ := redistore.NewRediStore(10, "tcp", ":6379", "", "", []byte("fg-key"))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	// 跨站请求需要 SameSite=None（配合 Secure=true）
	store.Options.SameSite = http.SameSiteNoneMode
	store.SetMaxAge(maxAge)

	gothic.Store = store
	providers := []goth.Provider{}

	// Google provider
	googleKey := getConfig("keys.GOOGLE_KEY", "GOOGLE_KEY")
	googleSecret := getConfig("keys.GOOGLE_SECRET", "GOOGLE_SECRET")
	if googleKey != "" && googleSecret != "" {
		accountGoogle := getConfig("keys.ACCOUNTS_GOOGLE", "ACCOUNTS_GOOGLE")
		oauth2Google := getConfig("keys.OAUTH2_GOOGLEAPIS", "OAUTH2_GOOGLEAPIS")
		openidconnectGoogle := getConfig("keys.OPENIDCONNECT_GOOGLEAPIS", "OPENIDCONNECT_GOOGLEAPIS")
		if accountGoogle != "" && oauth2Google != "" && openidconnectGoogle != "" {
			google.Endpoint = oauth2.Endpoint{
				AuthURL:       fmt.Sprintf("https://%s/o/oauth2/auth", accountGoogle),
				TokenURL:      fmt.Sprintf("https://%s/token", oauth2Google),
				DeviceAuthURL: fmt.Sprintf("https://%s/device/code", oauth2Google),
				AuthStyle:     oauth2.AuthStyleInParams,
			}
		}
		gg := google.New(googleKey, googleSecret, core.FillDomain("/auth/callback/google"))
		if openidconnectGoogle != "" {
			gg.EndpointProfile = fmt.Sprintf("https://%s/v1/userinfo", openidconnectGoogle)
		}
		providers = append(providers, gg)
		glog.InfoF("google provider inited")
	}

	// GitHub provider
	githubKey := getConfig("keys.GITHUB_KEY", "GITHUB_KEY")
	githubSecret := getConfig("keys.GITHUB_SECRET", "GITHUB_SECRET")
	if githubKey != "" && githubSecret != "" {
		githubDomain := getConfig("keys.GITHUB_DOMAIN", "GITHUB_DOMAIN")
		apiGithubDomain := getConfig("keys.API_GITHUB_DOMAIN", "API_GITHUB_DOMAIN")
		if githubDomain != "" && apiGithubDomain != "" {
			// 	AuthURL    = "https://github.com/login/oauth/authorize"
			//	TokenURL   = "https://github.com/login/oauth/access_token"
			//	ProfileURL = "https://api.github.com/user"
			//	EmailURL   = "https://api.github.com/user/emails"
			providers = append(providers, github.NewCustomisedURL(githubKey,
				githubSecret,
				core.FillDomain("/auth/callback/github"),
				fmt.Sprintf("https://%s/login/oauth/authorize", githubDomain),
				fmt.Sprintf("https://%s/login/oauth/access_token", githubDomain),
				fmt.Sprintf("https://%s/user", apiGithubDomain),
				fmt.Sprintf("https://%s/user/emails", apiGithubDomain),
			))
		} else {
			providers = append(providers, github.New(githubKey, githubSecret, core.FillDomain("/auth/callback/github")))
		}
		glog.InfoF("github provider inited")
	}

	// WeChat provider (微信开放平台 - PC扫码登录)
	wechatKey := getConfig("keys.WECHAT_KEY", "WECHAT_KEY")
	wechatSecret := getConfig("keys.WECHAT_SECRET", "WECHAT_SECRET")
	if wechatKey != "" && wechatSecret != "" {
		providers = append(providers, wechat.New(wechatKey, wechatSecret, core.FillDomain("/auth/callback/wechat"), wechat.WECHAT_LANG_CN))
		glog.InfoF("wechat provider inited")
	}

	if len(providers) == 0 {
		panic("no oauth provider configured")
	}

	goth.UseProviders(providers...)
}
