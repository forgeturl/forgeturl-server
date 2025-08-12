package middleware

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/lib/prometheus"
	"github.com/sunmi-OS/gocore/v2/utils"
)

func NewGin() *gin.Engine {
	if utils.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = NewGlogWriterDebug()
	gin.DefaultErrorWriter = NewGlogWriterError()
	r := gin.New()
	if utils.GetRunTime() != utils.LocalEnv {
		r.Use(gin.Recovery())
	}
	r.Use(ServerLogging(WithSlowThreshold(5000)))

	r.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome GoCore Service")
	})
	pprof.Register(r, "/debug/pprof")
	prometheus.NewPrometheus("app", nil).Use(r)
	return r
}

type GlogWriter struct {
	LogFunc func(args ...interface{})
}

func (w *GlogWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	w.LogFunc(string(p[:len(p)-1]))
	return len(p) - 1, nil
}

func NewGlogWriterDebug() *GlogWriter {
	return &GlogWriter{LogFunc: glog.Debug}
}
func NewGlogWriterError() *GlogWriter {
	return &GlogWriter{LogFunc: glog.Error}
}

func NewGlogWriterFatal() *GlogWriter {
	return &GlogWriter{LogFunc: glog.Fatal}
}
