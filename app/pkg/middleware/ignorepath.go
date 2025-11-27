package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IgnoreNotExistPath() gin.HandlerFunc {
	allowListMap := map[string]bool{
		"/health":     true,
		"/debug/":     true,
		"/monitor/":   true,
		"/login/":     true,
		"/dumplinks/": true,
		"/space/":     true,
	}
	// 如果不在前缀树里，则直接404，不记录到promethues里
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 提取path后内容直到第二个/
		firstIndex := strings.Index(path, "/")
		if firstIndex == -1 {
			c.Header("X-NotFound-Reason", "invalid-path")
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		secondIndex := strings.Index(path[firstIndex+1:], "/")
		var subPath string
		if secondIndex == -1 {
			// 没有第二个/，直接使用整个路径
			subPath = path
		} else {
			// 有第二个/，提取到第二个/为止（包含第二个/）
			subPath = path[:firstIndex+secondIndex+2]
		}

		if _, ok := allowListMap[subPath]; !ok {
			c.Header("X-NotFound-Reason", "not-in-allowlist")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()

	}
}

func MustLocalIp() gin.HandlerFunc {
	// 需要限制为本地IP访问的路径
	localOnlyPaths := map[string]bool{
		"/health":             true,
		"/debug/pprof":        true,
		"/monitor/prometheus": true,
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 提取path前缀，判断是否需要本地IP限制
		var needCheck bool
		for prefix := range localOnlyPaths {
			if strings.HasPrefix(path, prefix) {
				needCheck = true
				break
			}
		}

		// 如果不在限制列表中，直接放行
		if !needCheck {
			c.Next()
			return
		}

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 检查是否是本地或内网IP
		if !isLocalOrPrivateIP(clientIP) {
			c.Header("X-Forbidden-Reason", "only-local-ip-allowed")
			c.Header("X-Client-IP", clientIP) // 方便调试
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// isLocalOrPrivateIP 判断是否是本地或内网IP
func isLocalOrPrivateIP(ipStr string) bool {
	// 解析IP地址
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// 检查是否是回环地址 (127.0.0.1, ::1)
	if ip.IsLoopback() {
		return true
	}

	// 检查是否是内网地址
	if ip.IsPrivate() {
		return true
	}

	return false
}
