package api

import (
	"net/http"

	"github.com/LucienVen/golang-auth-service/internal/appcontext"
	"github.com/LucienVen/golang-auth-service/internal/controller"
	"github.com/LucienVen/golang-auth-service/internal/middleware"
	"github.com/LucienVen/golang-auth-service/internal/service"
	"github.com/LucienVen/golang-auth-service/pkg/auth"
	"github.com/gin-gonic/gin"
)

// Router 路由管理器
type Router struct {
	engine      *gin.Engine
	controllers *controller.Container
	validator   *auth.TokenValidator // 统一管理TokenValidator实例
}

// NewRouter 创建路由管理器
func NewRouter(appCtx *appcontext.AppContext) *Router {
	// 创建 Gin 引擎
	engine := gin.New()

	// 创建JWT服务和统一验证器
	jwtService := service.NewJWTService()
	validator := auth.NewTokenValidator(jwtService)

	// 创建控制器容器，注入验证器
	controllers := controller.NewContainerWithValidator(appCtx, validator)

	return &Router{
		engine:      engine,
		controllers: controllers,
		validator:   validator, // 保存验证器实例用于创建Middleware
	}
}

// RegisterRoutes 注册所有路由
func (r *Router) RegisterRoutes() {
	// 设置全局中间件
	r.setupMiddleware()

	// 设置所有路由组
	r.setupAuthRoutes()
	r.setupHealthRoutes()
	r.setupAPIRoutes()

	// 开发环境测试路由
	r.setupTestRoutes()
}

// setupMiddleware 设置全局中间件
func (r *Router) setupMiddleware() {
	// 1. 异常恢复 - 最外层，捕获所有panic
	r.engine.Use(gin.Recovery())

	// 2. 请求ID - 最优先设置，确保所有后续中间件都能使用
	r.engine.Use(middleware.RequestIDMiddleware())

	// 3. CORS处理 - 在业务逻辑之前处理跨域
	r.engine.Use(middleware.CORSMiddleware())

	// 4. 日志记录 - 最内层，能获取到完整的请求信息
	r.engine.Use(middleware.Logger())
}

// setupAuthRoutes 设置认证相关路由
func (r *Router) setupAuthRoutes() {
	// 创建认证中间件 - 使用Router中的统一验证器
	authMiddleware := middleware.NewAuthMiddleware(r.validator)

	// API版本组
	v1 := r.engine.Group("/api/v1")
	{
		// 认证路由（不需要认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.controllers.Auth.Login)
			auth.POST("/register", r.controllers.Auth.Register)
			auth.POST("/refresh", r.controllers.Auth.RefreshToken)
			auth.POST("/logout", authMiddleware.RequireAuth(), r.controllers.Auth.Logout)
			auth.GET("/validate", authMiddleware.RequireAuth(), r.controllers.Auth.ValidateToken)
		}

		// 用户路由（需要认证）
		user := v1.Group("/user")
		user.Use(authMiddleware.RequireAuth())
		{
			user.GET("/profile", r.controllers.Auth.GetProfile)
			user.GET("/sessions", r.controllers.Auth.GetSessions)
			user.DELETE("/sessions", r.controllers.Auth.RevokeAllSessions)
			user.DELETE("/sessions/:sessionId", r.controllers.Auth.RevokeSession)
		}
	}
}

// setupHealthRoutes 设置健康检查路由
func (r *Router) setupHealthRoutes() {
	// 健康检查路由（无需认证）
	r.engine.GET("/health", r.controllers.Health.Check)

	// API信息路由
	r.engine.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "golang-auth-service",
			"version":     "1.0.0",
			"description": "Go语言认证服务",
		})
	})

	// Ping路由（兼容性）
	r.engine.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

// setupAPIRoutes 设置其他API路由（保持兼容性）
func (r *Router) setupAPIRoutes() {
	// 基础路由组（保持兼容性）
	base := r.engine.Group("/api")
	{
		// 健康检查（重复路由，但保持兼容性）
		base.GET("/health", r.controllers.Health.Check)
	}
}

// setupTestRoutes 设置测试路由（仅用于开发环境）
func (r *Router) setupTestRoutes() {
	// 创建认证中间件 - 使用Router中的统一验证器
	authMiddleware := middleware.NewAuthMiddleware(r.validator)

	test := r.engine.Group("/test")
	{
		test.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		test.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
			userID, _ := middleware.GetCurrentUserID(c)
			username, _ := middleware.GetCurrentUsername(c)
			c.JSON(200, gin.H{
				"message":  "访问受保护的资源成功",
				"user_id":  userID,
				"username": username,
			})
		})

		test.GET("/optional", authMiddleware.OptionalAuth(), func(c *gin.Context) {
			userID, exists := middleware.GetCurrentUserID(c)
			if exists {
				username, _ := middleware.GetCurrentUsername(c)
				c.JSON(200, gin.H{
					"message":       "可选认证成功",
					"user_id":       userID,
					"username":      username,
					"authenticated": true,
				})
			} else {
				c.JSON(200, gin.H{
					"message":       "未认证访问",
					"authenticated": false,
				})
			}
		})
	}
}

// GetEngine 获取 Gin 引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
