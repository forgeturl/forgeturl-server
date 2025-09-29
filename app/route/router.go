package route

import (
	"forgeturl-server/api"
	"forgeturl-server/api/dumplinks"
	"forgeturl-server/api/login"
	"forgeturl-server/api/space"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	loginService := api.NewLoginService()
	space.RegisterSpaceServiceHTTPServer(router, api.NewSpaceService())
	login.RegisterLoginServiceHTTPServer(router, loginService)
	dumplinks.RegisterDumplinksServiceHTTPServer(router, api.NewDumplinksService())

	router.GET("/login/connector/auth", api.LoginAuth(loginService))                   // 连接器登录，跳转鉴权的url
	router.GET("/login/connector/callback/:provider", api.LoginCallback(loginService)) // 第三方登录回调
}
