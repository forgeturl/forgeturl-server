package route

import (
	"2049links-server/api"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	user := router.Group("/app/user")
	user.POST("/getUserInfo", api.GetUserInfo) //获取用户信息

}
