package controller

import (
	"github.com/LucienVen/golang-auth-service/internal/appcontext"
)

// Container 控制器容器
type Container struct {
	Health *HealthController
	// 在这里添加其他控制器
}

// NewContainer 创建控制器容器
func NewContainer(appCtx *appcontext.AppContext) *Container {
	return &Container{
		Health: NewHealthController(appCtx.DB),
		// 在这里初始化其他控制器
	}
}
