package controller

import (
	"github.com/LucienVen/golang-auth-service/internal/db"
	"github.com/LucienVen/golang-auth-service/internal/errors"
	"github.com/LucienVen/golang-auth-service/internal/response"
	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct {
	db db.DB
}

// NewHealthController 创建健康检查控制器
func NewHealthController(db db.DB) *HealthController {
	return &HealthController{
		db: db,
	}
}

// Check 执行健康检查
func (c *HealthController) Check(ctx *gin.Context) {
	if err := c.db.Ping(); err != nil {
		response.ServerError(ctx, err)
		return
	}
	response.Success(ctx, gin.H{
		"status": "ok",
		"code":   errors.ErrCodeSuccess,
	})
}
