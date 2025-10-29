package controller

import (
	"github.com/LucienVen/golang-auth-service/internal/appcontext"
	"github.com/LucienVen/golang-auth-service/internal/service"
)

// Container 控制器容器
type Container struct {
	Health *HealthController
	Auth   *AuthController // 添加认证控制器
}

// NewContainer 创建控制器容器
func NewContainer(appCtx *appcontext.AppContext) *Container {
	// 创建认证服务
	authService := service.NewAuthService(appCtx)

	return &Container{
		Health: NewHealthController(appCtx.DB),
		Auth:   NewAuthController(authService),
	}
}
