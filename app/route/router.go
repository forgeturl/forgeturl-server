package route

import (
	"forgeturl-server/api"
	"forgeturl-server/api/space"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	space.RegisterSpaceServiceHTTPServer(router, api.NewSpaceService())
}
