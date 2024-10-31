package api

import "github.com/gin-gonic/gin"

func getLoginState(ginCtx *gin.Context) (isLogin bool, uid int64) {
	auth := ginCtx.GetHeader("Authorization")
	if auth == "" {
		return false, 0
	}

	// 通过缓存获取登录账号信息
	return
}
