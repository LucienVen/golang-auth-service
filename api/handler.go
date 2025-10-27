package api

import (
	"net/http"

	"github.com/LucienVen/golang-auth-service/internal/appcontext"
	"github.com/LucienVen/golang-auth-service/internal/controller"
	"github.com/gin-gonic/gin"
)

// Handler 处理函数结构体
type Handler struct {
	controllers *controller.Container
}

// NewHandler 创建处理函数
func NewHandler(appCtx *appcontext.AppContext) *Handler {
	return &Handler{
		controllers: controller.NewContainer(appCtx),
	}
}

// HealthCheck 健康检查处理函数
func (h *Handler) HealthCheck(c *gin.Context) {
	h.controllers.Health.Check(c)
}

// Ping 测试连接处理函数
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
