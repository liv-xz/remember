package middleware

import (
	"sync/atomic"
	"time"

	"server/applog"

	"github.com/gin-gonic/gin"
)

var totalRequests int64

func GetTotalRequests() int64 {
	return atomic.LoadInt64(&totalRequests)
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		n := atomic.AddInt64(&totalRequests, 1)
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		if status >= 400 {
			applog.GetRequestLog(start, c.Request.Method, c.Request.URL.Path, status, duration)
		} else if n%20 == 0 {
			applog.GetRequestLog(start, c.Request.Method, c.Request.URL.Path, status, duration)
		}
	}
}
