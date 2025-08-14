package middleware

import (
	"forgeturl-server/dal"

	"github.com/sunmi-OS/gocore/v2/api"
)

// GetLoginUid 获取登录者的uid
// 一般是api的第一个接口
func GetLoginUid(g *api.Context) int64 {
	token := g.GetHeader("X-Token")
	// x-token:{} -> userid
	// 这里可以通过内存缓存一份
	return dal.C.GetXToken(g.Request.Context(), token)
}
