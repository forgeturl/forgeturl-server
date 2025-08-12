package middleware

import (
	"bytes"
	"io"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"
)

const hideBody = "gocore_body_hide"
const maxBodySize = 1024 * 1024 * 1

// 自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	// 向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	// 完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

// ServerLogging middleware for accesslog
func ServerLogging(options ...Option) gin.HandlerFunc {
	op := option{
		slowThresholdMs:     1000,
		hideLogsWithPath:    hideLogsPath,
		hideReqBodyWithPath: map[string]bool{},
		hideRespBodWithPath: hidelRespBodyLogsPath,
		allowShowHeaders:    map[string]bool{},
		hideShowHeaders:     hideShowHeaders,
	}
	for _, apply := range options {
		apply(&op)
	}

	return func(c *gin.Context) {
		r := c.Request
		path := r.URL.Path
		start := time.Now()
		quota := int64(-1)
		if deadline, ok := r.Context().Deadline(); ok {
			quota = time.Until(deadline).Milliseconds()
		}
		body := ""

		if op.hideReqBodyWithPath[path] || c.Request.ContentLength > maxBodySize {
			body = hideBody
		} else {
			b, err := c.GetRawData()
			if err != nil {
				body = "failed to get request body"
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				body = string(b)
			}
		}

		hideResp := op.hideRespBodWithPath[path]
		var writer responseWriter
		if !hideResp {
			writer = responseWriter{
				c.Writer,
				bytes.NewBuffer([]byte{}),
			}
			c.Writer = writer
		}
		traceid := c.GetHeader(utils.XB3TraceId)
		if traceid == "" { // 如果找不到x-b3-traceid，用x-request-id
			traceid = c.GetHeader(utils.XRequestId)
		}
		if traceid == "" { // 还找不到，自己生成一个大写的traceid
			traceid = NewUUID()
			traceid = strings.ToUpper(traceid)
		}
		ctx := utils.SetMetaData(c.Request.Context(), utils.XB3TraceId, traceid)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-AppName", utils.GetAppName())

		c.Next()

		r = c.Request
		ctx = r.Context()
		responseCode := math.MinInt8
		var responseMsg string
		var respBytes []byte
		if !hideResp {
			respBytes = writer.b.Bytes()
			if root, err0 := sonic.Get(respBytes); err0 == nil {
				code, err := root.Get("code").Int64()
				if err == nil {
					responseCode = int(code)
				}
				msg, err := root.Get("msg").String()
				if err == nil {
					responseMsg = msg
				}
			}

			if len(respBytes) > maxBodySize {
				respBytes = []byte(hideBody)
			}
		} else {
			respBytes = []byte(hideBody)
		}

		reqAppname := r.Header.Get(utils.XAppName)
		statusCode := c.Writer.Status()
		costms := time.Since(start).Milliseconds()

		if op.hideLogsWithPath[path] {
			return
		}

		remoteIp, remotePort, _ := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
		fields := []interface{}{
			"kind", "server",
			"costms", costms,
			"client_ip", c.ClientIP(),
			"remote_ip", remoteIp,
			"remote_port", remotePort,
			"host", r.Host,
			"method", r.Method,
			"path", path,
			"req", utils.LogContentUnmarshal(body),
			"resp", utils.LogContentUnmarshal(string(respBytes)),
			"status", statusCode, // http状态码
			"code", responseCode, // 业务错误码
			"msg", responseMsg,

			"start_time", start.Format(utils.TimeFormat),
			"req_header", filterHeaders(r.Header, op.allowShowHeaders, op.hideShowHeaders),
			"resp_header", c.Writer.Header(),
		}
		if reqAppname != "" {
			fields = append(fields, "req_appname", reqAppname)
		}
		if r.URL.RawQuery != "" {
			fields = append(fields, "params", r.URL.RawQuery)
		}
		if c.GetHeader("x-forwarded-for") != "" {
			fields = append(fields, "forward_ip", c.GetHeader("x-forwarded-for"))
		}
		if quota != -1 {
			fields = append(fields, "timeout_quota", quota) // 收到请求时，剩余处理时间
		}

		logFunc := glog.InfoV
		isErrorLevel := false
		if statusCode >= http.StatusInternalServerError {
			logFunc = glog.ErrorV
			isErrorLevel = true
		} else if statusCode >= http.StatusBadRequest {
			logFunc = glog.WarnV
		} else if op.slowThresholdMs != 0 && costms > op.slowThresholdMs {
			logFunc = glog.WarnV
		}

		if responseCode >= 19105000 && responseCode < 19106000 {
			logFunc = glog.ErrorV
		} else if !isErrorLevel && responseCode >= 19104000 && responseCode < 19105000 {
			logFunc = glog.WarnV
		}

		logFunc(ctx, fields...)
	}
}

func mustPositive(val float64) float64 {
	if val < 0 {
		return 0
	}
	return val
}

const maxShowHeaderLen = 32

// header白名单过滤
func filterHeaders(headers http.Header, allowShowHeaders map[string]bool, hideShowHeaders map[string]bool) http.Header {
	filteredHeaders := http.Header{}
	for k, v := range headers {
		if len(filteredHeaders) >= maxShowHeaderLen {
			break
		}
		lower := strings.ToLower(k)
		// 优先判断白名单
		if allowShowHeaders[k] {
			filteredHeaders[k] = v
			continue
		}
		if hideShowHeaders[lower] {
			continue
		}
		// 如果没配置白名单，能走到这里则允许打印
		if len(allowShowHeaders) == 0 {
			filteredHeaders[k] = v
		}
	}
	return filteredHeaders
}

// 黑名单 某些路径不打印response body，但打印日志
var hidelRespBodyLogsPath = map[string]bool{
	"/private/debug/pprof/profile": true,
	"/debug/pprof/":                true,
	"/debug/pprof/cmdline":         true,
	"/debug/pprof/profile":         true,
	"/debug/pprof/symbol":          true,
	"/debug/pprof/trace":           true,
	"/monitor/prometheus":          true,
}

// 黑名单 某些路径不打印日志
var hideLogsPath = map[string]bool{
	"/health": true,
}

var hideShowHeaders = map[string]bool{
	"accept":                true,
	"accept-encoding":       true,
	"proxy-connection":      true,
	"x-envoy-peer-metadata": true,
}

// WithSlowThreshold 当请求耗时超过slowThreshold时，打印slow log。建议配置1000
func WithSlowThreshold(slowThresholdMs int64) Option {
	return func(o *option) {
		o.slowThresholdMs = slowThresholdMs
	}
}

// WithHideLogsPath 对某些路径不打印日志
func WithHideLogsPath(hideLogsWithPath map[string]bool, isMerge bool) Option {
	return func(o *option) {
		if isMerge {
			o.hideLogsWithPath = mergeMap(o.hideLogsWithPath, hideLogsWithPath)
		} else {
			o.hideLogsWithPath = hideLogsWithPath
		}
	}
}

// WithHideBodyLogsPath 对某些路径不打印body
func WithHideBodyLogsPath(hideBodyLogsWithPath map[string]bool, isMerge bool) Option {
	return func(o *option) {
		if isMerge {
			o.hideRespBodWithPath = mergeMap(o.hideRespBodWithPath, hideBodyLogsWithPath)
		} else {
			o.hideRespBodWithPath = hideBodyLogsWithPath
		}
	}
}

func WithHideReqBodyLogsPath(hideBodyLogsWithPath map[string]bool, isMerge bool) Option {
	return func(o *option) {
		if isMerge {
			o.hideReqBodyWithPath = mergeMap(o.hideReqBodyWithPath, hideBodyLogsWithPath)
		} else {
			o.hideReqBodyWithPath = hideBodyLogsWithPath
		}
	}
}

// WithAllowShowHeaders 只展示某些header
func WithAllowShowHeaders(allowHeaders []string) Option {
	return func(o *option) {
		for _, header := range allowHeaders {
			o.allowShowHeaders[strings.ToLower(header)] = true
		}
	}
}

func WithHideShowHeaders(hideHeaders map[string]bool, isMerge bool) Option {
	return func(o *option) {
		if isMerge {
			o.hideShowHeaders = mergeMap(o.hideShowHeaders, hideHeaders)
		} else {
			o.hideShowHeaders = hideHeaders
		}
	}
}

func mergeMap(m1, m2 map[string]bool) map[string]bool {
	for k, v := range m2 {
		if _, ok := m1[k]; !ok {
			m1[k] = v
		}
	}
	return m1
}

type option struct {
	slowThresholdMs     int64
	hideLogsWithPath    map[string]bool
	hideReqBodyWithPath map[string]bool
	hideRespBodWithPath map[string]bool
	allowShowHeaders    map[string]bool
	hideShowHeaders     map[string]bool
}

type Option func(op *option)

func NewUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
