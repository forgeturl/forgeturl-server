package middleware

import (
	"forgeturl-server/dal"
	"strings"

	"github.com/sunmi-OS/gocore/v2/api"
)

// GetLoginUid 获取登录者的uid
// 一般是api的第一个接口
func GetLoginUid(g *api.Context) int64 {
	ctx := g.Request.Context()
	if strings.HasPrefix(g.Request.URL.Path, "/openclaw/") {
		authHeader := strings.TrimSpace(g.GetHeader("Authorization"))
		if len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
			apiKey := strings.TrimSpace(authHeader[7:])
			if apiKey != "" {
				if uid := dal.OpenClaw.GetUIDByAPIKey(ctx, apiKey); uid > 0 {
					return uid
				}
			}
		}
	}

	token := g.GetHeader("X-Token")
	// x-token:{} -> userid
	// 这里可以通过内存缓存一份
	return dal.C.GetXToken(ctx, token)
}
