package controller

import (
	"github.com/LucienVen/golang-auth-service/internal/db"
)

// HealthController 健康检查控制器
type HealthController struct {
	db *db.GormDB
}

// NewHealthController 创建健康检查控制器
func NewHealthController(db *db.GormDB) *HealthController {
	return &HealthController{
		db: db,
	}
}

// Check 执行健康检查
func (c *HealthController) Check() error {
	return c.db.Ping()
}
