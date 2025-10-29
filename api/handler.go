package api

import (
	"github.com/LucienVen/golang-auth-service/internal/appcontext"
	"github.com/LucienVen/golang-auth-service/internal/controller"
)

// Handler 处理函数结构体 - 仅作为适配器保持兼容性
type Handler struct {
	controllers *controller.Container
}

// NewHandler 创建处理函数 - 已弃用，建议直接使用 api.Router
// @deprecated: 推荐使用 NewRouter 代替
func NewHandler(appCtx *appcontext.AppContext) *Handler {
	return &Handler{
		controllers: controller.NewContainer(appCtx),
	}
}

// 所有处理函数已迁移到 api/router.go 中
// HealthCheck、Ping 等功能现在通过控制器统一管理
