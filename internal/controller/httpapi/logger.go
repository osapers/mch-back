package httpapi

import (
	"time"

	"github.com/osapers/mch-back/internal/constant"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type requestInfo struct {
	Status           int
	Method, Path, IP string
	Time             time.Time
	Latency          time.Duration
}

func (r *requestInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("status", r.Status)
	enc.AddString("method", r.Method)
	enc.AddString("path", r.Path)
	enc.AddString("started", r.Time.Format(time.RFC3339))
	enc.AddString("latency", r.Latency.String())
	return nil
}

func loggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		reqInfo := zap.Object("reqInfo", &requestInfo{
			Status:  c.Writer.Status(),
			Method:  method,
			Path:    path,
			IP:      c.ClientIP(),
			Time:    end,
			Latency: latency,
		})

		if len(c.Errors) > 0 {
			logger.With(reqInfo).Error("Error response", zap.Error(c.Err()))
		}

		if err, exists := c.Get(constant.ErrorCtxKey.String()); exists && err != nil {
			logger.With(reqInfo).Error("Error response", zap.Error(err.(error)))
		} else {
			logger.With(reqInfo).Info("Success response")
		}

		c.Next()
	}
}
