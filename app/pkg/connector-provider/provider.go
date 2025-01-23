package connector_provider

import (
	"github.com/markbates/goth/providers/github"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/wechat"
)

const (
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

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "https://api.forgeturl.com/login/connector/google"),
		wechat.New(os.Getenv("WECHAT_KEY"), os.Getenv("WECHAT_SECRET"), "https://api.forgeturl.com/login/connector/wechat", wechat.WECHAT_LANG_CN),
		//github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "https://api.forgeturl.com/login/connector/github"),
		github.New(clientKey, clientSecret, "http://localhost:8080/auth/github/callback"),
	)

	g, err := goth.GetProvider("google")
	if err != nil {
		panic(err)
	}
	_ = g

	gP := google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "https://api.forgeturl.com/login/connector/google")
	gP.SetHostedDomain()

}
