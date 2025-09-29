package connector_provider

import (
	"os"

	"forgeturl-server/pkg/core"

	"github.com/forgeturl/redistore"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/wechat"
)

func Init() {
	// todo 自己实现一个，通过header返回， 不通过cookie返回

	maxAge := 86400 * 30 // 30 days
	isProd := false
	store, _ := redistore.NewRediStore(10, "tcp", ":6379", "", "", []byte("fg-key"))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.SetMaxAge(maxAge)

	gothic.Store = store
	viders := []goth.Provider{
		//facebook.New(),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), core.FillDomain("/login/connector/callback/google")),
		wechat.New(os.Getenv("WECHAT_KEY"), os.Getenv("WECHAT_SECRET"), core.FillDomain("/login/connector/callback/wechat"), wechat.WECHAT_LANG_CN),
	}
	// Add GitHub provider if environment variables are set
	if os.Getenv("GITHUB_KEY") != "" && os.Getenv("GITHUB_SECRET") != "" {
		viders = append(viders, github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), core.FillDomain("/login/connector/callback/github")))
	}
	goth.UseProviders(viders...)

	_, err := goth.GetProvider("google")
	if err != nil {
		panic(err)
	}
}
