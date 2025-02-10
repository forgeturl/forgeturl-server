package connector_provider

import (
	"forgeturl-server/pkg/core"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/wechat"
)

const (
	// 使用goth例子中的clientKey和clientSecret，仅用于测试，非自己申请
	clientKey    = "dfc0084166d0ef8a9aac"
	clientSecret = "75a7af7446b893707299595d1a4718b7c81174a8"
)

func Init() {
	// store, _ = redistore.NewRediStore(10, "tcp", ":6379", "", []byte("redis-key"))
	key := ""
	maxAge := 86400 * 30 // 30 days
	isProd := false

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store
	viders := []goth.Provider{
		//facebook.New(),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), core.FillDomain("/login/connector/google")),
		wechat.New(os.Getenv("WECHAT_KEY"), os.Getenv("WECHAT_SECRET"), core.FillDomain("/login/connector/wechat"), wechat.WECHAT_LANG_CN),
	}
	if core.IsTestEnv() {
		viders = append(viders, github.New(clientKey, clientSecret, core.FillDomain("/login/connector/callback/")))
	}
	goth.UseProviders(viders...)

	g, err := goth.GetProvider("google")
	if err != nil {
		panic(err)
	}
	_ = g

}
