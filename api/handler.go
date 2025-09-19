package api

import (
	"net/http"

	"github.com/LucienVen/golang-auth-service/internal/controller"
	"github.com/LucienVen/golang-auth-service/internal/db"
	"github.com/gin-gonic/gin"
)

// Handler 处理函数结构体
type Handler struct {
	controllers *controller.Container
}

// NewHandler 创建处理函数
func NewHandler(db *db.GormDB) *Handler {
	return &Handler{
		controllers: controller.NewContainer(db),
	}
}

// HealthCheck 健康检查处理函数
func (h *Handler) HealthCheck(c *gin.Context) {
	if err := h.controllers.Health.Check(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Ping 测试连接处理函数
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
