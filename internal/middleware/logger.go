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

		// 获取请求ID - 安全处理nil值
		var requestIDStr string
		requestID, _ := c.Get("request_id")
		if requestID != nil {
			if id, ok := requestID.(string); ok {
				requestIDStr = id
			} else {
				requestIDStr = ""
			}
		} else {
			requestIDStr = ""
		}

		// 获取用户信息（如果已认证）
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 构建日志字段
		fields := []zap.Field{
			zap.String("request_id", requestIDStr),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.String("ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.GetHeader("User-Agent")),
		}

		// 如果用户已认证，添加用户信息
		if userID != nil {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}
		if username != nil {
			fields = append(fields, zap.String("username", username.(string)))
		}

		// 记录日志
		log.Info("HTTP Request", fields...)
	}
}
