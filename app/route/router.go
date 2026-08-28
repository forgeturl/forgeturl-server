package route

import (
	"forgeturl-server/api"
	"forgeturl-server/api/dumplinks"
	"forgeturl-server/api/login"
	"forgeturl-server/api/openclaw"
	"forgeturl-server/api/space"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/lib/prometheus"
)

func Routes(router *gin.Engine) {
	space.RegisterSpaceServiceHTTPServer(router, api.NewSpaceService())
	login.RegisterLoginServiceHTTPServer(router, api.NewLoginService())
	dumplinks.RegisterDumplinksServiceHTTPServer(router, api.NewDumplinksService())
	openclaw.RegisterOpenClawServiceHTTPServer(router, api.NewOpenClawService())
	space.SetAutoValidate(false, nil, true)
	login.SetAutoValidate(false, nil, true)
	dumplinks.SetAutoValidate(false, nil, true)
	openclaw.SetAutoValidate(false, nil, true)

	router.GET("/login/connector/auth", api.LoginAuth())                   // 连接器登录，跳转鉴权的url
	router.GET("/login/connector/callback/:provider", api.LoginCallback()) // 第三方登录回调
	router.POST("/login/connector/avm/exchange", api.AVMAuthCodeExchange())
	router.GET("/login/connector/avm/wechat-mp/bind", api.AVMWeChatMPBind())
	router.GET("/login/connector/avm/wechat-mp/callback", api.AVMWeChatMPCallback())
	router.POST("/login/connector/avm/wechat-mp/exchange", api.AVMWeChatMPBindExchange())
	router.POST("/login/connector/avm/wechat-mp/send", api.AVMWeChatMPSend())

	router.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome GoCore Service")
	})

	// Stats API - 公开接口，获取用户统计信息
	router.GET("/stats/users", api.GetUserCount())
	pprof.Register(router, "/debug/pprof")
	prometheus.NewPrometheus("app", nil).Use(router)
}
