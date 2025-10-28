package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/LucienVen/golang-auth-service/internal/appcontext"
	"github.com/LucienVen/golang-auth-service/internal/router"
	"github.com/LucienVen/golang-auth-service/internal/service"
)

func main() {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建应用上下文（简化版本）
	appCtx := &appcontext.AppContext{}

	// 创建认证服务
	authService := service.NewAuthService(appCtx)

	// 创建路由
	authRoutes := router.NewAuthRoutes(authService)

	// 创建Gin引擎
	router := gin.New()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 设置路由
	authRoutes.SetupRoutes(router)

	// 在开发环境添加测试路由
	authRoutes.SetupTestRoutes(router)

	// 启动服务器
	port := ":8080"
	log.Printf("认证服务启动，监听端口 %s", port)
	log.Printf("API文档: http://localhost%s/api", port)
	log.Printf("健康检查: http://localhost%s/health", port)
	log.Printf("测试接口: http://localhost%s/test/ping", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}