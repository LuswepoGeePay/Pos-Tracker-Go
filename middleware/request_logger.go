package middleware

import (
	"pos-master/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)

		attrs := []any{
			"method", method,
			"path", path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"client_ip", clientIP,
		}

		if errMsg, ok := c.Get(utils.ResponseErrorKey); ok {
			attrs = append(attrs, "error", errMsg)
		}
		if detail, ok := c.Get(utils.ResponseErrorDetailKey); ok {
			attrs = append(attrs, "detail", detail)
		}

		if status >= 400 {
			utils.Warn("request failed", attrs...)
			return
		}

		utils.Info("request completed", attrs...)
	}
}
