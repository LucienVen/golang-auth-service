package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/LucienVen/golang-auth-service/internal/errors"
	"github.com/LucienVen/golang-auth-service/internal/response"
	"github.com/LucienVen/golang-auth-service/internal/service"
)

// AuthMiddleware JWT认证中间件
type AuthMiddleware struct {
	authService service.AuthService
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.ErrCodeUnauthorized, "缺少认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Error(c, errors.ErrCodeUnauthorized, "无效的令牌格式")
			c.Abort()
			return
		}

		// 提取令牌
		tokenString := authHeader[len(bearerPrefix):]
		if tokenString == "" {
			response.Error(c, errors.ErrCodeUnauthorized, "令牌不能为空")
			c.Abort()
			return
		}

		// 验证令牌
		userResponse, err := m.authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			response.Error(c, errors.ErrCodeUnauthorized, "无效的令牌")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", userResponse)
		c.Set("user_id", userResponse.ID)
		c.Set("username", userResponse.Username)

		c.Next()
	}
}

// OptionalAuth 可选认证的中间件（不强制要求认证）
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Bearer前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.Next()
			return
		}

		// 提取令牌
		tokenString := authHeader[len(bearerPrefix):]
		if tokenString == "" {
			c.Next()
			return
		}

		// 验证令牌
		userResponse, err := m.authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", userResponse)
		c.Set("user_id", userResponse.ID)
		c.Set("username", userResponse.Username)

		c.Next()
	}
}

// RequireRole 需要特定角色的中间件（扩展功能）
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先确保用户已认证
		user, exists := c.Get("user")
		if !exists {
			response.Error(c, errors.ErrCodeUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		userResponse, ok := user.(*response.UserResponse)
		if !ok {
			response.Error(c, errors.ErrCodeUnauthorized, "无效的用户信息")
			c.Abort()
			return
		}

		// TODO: 实现角色检查逻辑
		// 当前简化实现，直接通过
		_ = userResponse
		_ = roles

		c.Next()
	}
}

// GetCurrentUser 从上下文中获取当前用户
func GetCurrentUser(c *gin.Context) (*response.UserResponse, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userResponse, ok := user.(*response.UserResponse)
	return userResponse, ok
}

// GetCurrentUserID 从上下文中获取当前用户ID
func GetCurrentUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetCurrentUsername 从上下文中获取当前用户名
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	name, ok := username.(string)
	return name, ok
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简单实现，生产环境可以使用更复杂的算法
	return "req-" + strings.ReplaceAll(generateUUID(), "-", "")
}

// generateUUID 生成UUID v4
func generateUUID() string {
	// 使用Google UUID库生成真正的UUID v4
	return uuid.New().String()
}