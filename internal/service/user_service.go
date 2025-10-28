package service

import (
	"github.com/LucienVen/golang-auth-service/internal/appcontext"
)

// UserService 用户服务接口（简化版本）
type UserService interface {
	// 简化的用户服务接口，用于认证服务集成
}

// UserServiceImpl 用户服务实现（简化版本）
type UserServiceImpl struct {
	// 简化实现，不依赖外部存储
}

// NewUserService 创建用户服务实例
func NewUserService(appCtx *appcontext.AppContext) UserService {
	return &UserServiceImpl{}
}