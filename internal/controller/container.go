package controller

import (
	"github.com/LucienVen/golang-auth-service/internal/appcontext"
	"github.com/LucienVen/golang-auth-service/internal/service"
	"github.com/LucienVen/golang-auth-service/pkg/auth"
)

// Container 控制器容器
type Container struct {
	Health *HealthController
	Auth   *AuthController
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

// NewContainerWithValidator 使用外部验证器创建控制器容器
func NewContainerWithValidator(appCtx *appcontext.AppContext, validator *auth.TokenValidator) *Container {
	// 使用外部验证器创建认证服务
	authService := service.NewAuthServiceWithValidator(validator, appCtx)

	return &Container{
		Health: NewHealthController(appCtx.DB),
		Auth:   NewAuthController(authService),
	}
}
