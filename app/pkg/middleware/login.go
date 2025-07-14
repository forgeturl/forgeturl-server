package middleware

import "github.com/sunmi-OS/gocore/v2/api"

// GetLoginUid 获取登录者的uid
// 一般是api的第一个接口
func GetLoginUid(g *api.Context) int64 {
	token := g.GetHeader("X-Token")
	if token == "" {
		return 0
	}
	// todo 读取 redis，拉取到用户信息，并且设置到 gin.Context 中
	// x-token:{} -> userid
	// 这里可以通过内存缓存一份

	return 0
}
