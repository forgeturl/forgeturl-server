package route

import (
	"forgeturl-server/api"
	"forgeturl-server/api/dumplinks"
	"forgeturl-server/api/login"
	"forgeturl-server/api/space"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	space.RegisterSpaceServiceHTTPServer(router, api.NewSpaceService())
	login.RegisterLoginServiceHTTPServer(router, api.NewLoginService())
	dumplinks.RegisterDumplinksServiceHTTPServer(router, api.NewDumplinksService())

}
