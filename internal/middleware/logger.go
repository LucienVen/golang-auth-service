package middleware

import (
	"time"

	"github.com/LucienVen/golang-auth-service/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 自定义日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取状态码
		status := c.Writer.Status()

		// 记录日志
		log.Info("Gin Request",
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.String("ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.Duration("latency", latency),
		)
	}
}
