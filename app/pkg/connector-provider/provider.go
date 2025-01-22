package connector_provider

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/wechat"
)

func Init() {
	store, _ = redistore.NewRediStore(10, "tcp", ":6379", "", []byte("redis-key"))

	key := ""
	maxAge := 86400 * 30 // 30 days
	isProd := false

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "https://api.forgeturl.com/login/connector/google"),
		wechat.New(os.Getenv("WECHAT_KEY"), os.Getenv("WECHAT_SECRET"), "https://api.forgeturl.com/login/connector/wechat", wechat.WECHAT_LANG_CN),
	)

	g, err := goth.GetProvider("google")
	if err != nil {
		panic(err)
	}
	_ = g

	gP := google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "https://api.forgeturl.com/login/connector/google")
	gP.SetHostedDomain()

}
