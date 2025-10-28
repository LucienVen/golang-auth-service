package router

import (
	"github.com/gin-gonic/gin"
	"github.com/LucienVen/golang-auth-service/internal/controller"
	"github.com/LucienVen/golang-auth-service/internal/middleware"
	"github.com/LucienVen/golang-auth-service/internal/service"
)

// AuthRoutes 认证路由配置
type AuthRoutes struct {
	authController *controller.AuthController
	authMiddleware *middleware.AuthMiddleware
}

// NewAuthRoutes 创建认证路由
func NewAuthRoutes(authService service.AuthService) *AuthRoutes {
	return &AuthRoutes{
		authController: controller.NewAuthController(authService),
		authMiddleware: middleware.NewAuthMiddleware(authService),
	}
}

// SetupRoutes 设置认证相关路由
func (r *AuthRoutes) SetupRoutes(router *gin.Engine) {
	// 添加全局中间件
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())

	// API版本组
	v1 := router.Group("/api/v1")
	{
		// 认证路由（不需要认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authController.Login)
			auth.POST("/register", r.authController.Register)
			auth.POST("/refresh", r.authController.RefreshToken)
			auth.POST("/logout", r.authMiddleware.RequireAuth(), r.authController.Logout)
			auth.GET("/validate", r.authMiddleware.RequireAuth(), r.authController.ValidateToken)
		}

		// 用户路由（需要认证）
		user := v1.Group("/user")
		user.Use(r.authMiddleware.RequireAuth())
		{
			user.GET("/profile", r.authController.GetProfile)
			user.GET("/sessions", r.authController.GetSessions)
			user.DELETE("/sessions", r.authController.RevokeAllSessions)
			user.DELETE("/sessions/:sessionId", r.authController.RevokeSession)
		}
	}

	// 健康检查路由（无需认证）
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// API信息路由
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "golang-auth-service",
			"version":     "1.0.0",
			"description": "Go语言认证服务",
		})
	})
}

// SetupTestRoutes 设置测试路由（仅用于开发环境）
func (r *AuthRoutes) SetupTestRoutes(router *gin.Engine) {
	test := router.Group("/test")
	{
		test.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		test.GET("/protected", r.authMiddleware.RequireAuth(), func(c *gin.Context) {
			userID, _ := middleware.GetCurrentUserID(c)
			username, _ := middleware.GetCurrentUsername(c)
			c.JSON(200, gin.H{
				"message":  "访问受保护的资源成功",
				"user_id":  userID,
				"username": username,
			})
		})

		test.GET("/optional", r.authMiddleware.OptionalAuth(), func(c *gin.Context) {
			userID, exists := middleware.GetCurrentUserID(c)
			if exists {
				username, _ := middleware.GetCurrentUsername(c)
				c.JSON(200, gin.H{
					"message":  "可选认证成功",
					"user_id":  userID,
					"username": username,
					"authenticated": true,
				})
			} else {
				c.JSON(200, gin.H{
					"message": "未认证访问",
					"authenticated": false,
				})
			}
		})
	}
}