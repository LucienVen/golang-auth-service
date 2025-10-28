package service

import (
	"context"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/request"
	"github.com/LucienVen/golang-auth-service/internal/response"
)

// UserService 用户服务接口
type UserService interface {
	// Register 用户注册
	Register(ctx context.Context, req *request.RegisterRequest) (*response.RegisterResponse, error)

	// UpdateProfile 更新用户资料
	UpdateProfile(ctx context.Context, userID uint, req *request.UpdateProfileRequest) (*response.UpdateProfileResponse, error)

	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) (*response.ChangePasswordResponse, error)

	// GetUserProfile 获取用户信息
	GetUserProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error)

	// GetUserProfileByAccount 通过账号获取用户信息（内部使用）
	GetUserProfileByAccount(ctx context.Context, account string) (*response.UserProfileResponse, error)

	// ValidateAccountExists 检查账号是否存在
	ValidateAccountExists(ctx context.Context, req *request.RegisterRequest) error

	// ValidateUpdateAccountAvailable 验证更新账号的可用性
	ValidateUpdateAccountAvailable(ctx context.Context, userID uint, req *request.UpdateProfileRequest) error

	// DeactivateUser 停用用户（软删除）
	DeactivateUser(ctx context.Context, userID uint) error

	// ActivateUser 激活用户（管理员功能）
	ActivateUser(ctx context.Context, userID uint) error

	// DisableUser 禁用用户（管理员功能）
	DisableUser(ctx context.Context, userID uint) error

	// IsUserActive 检查用户是否处于活跃状态
	IsUserActive(ctx context.Context, userID uint) (bool, error)

	// ValidateUserAccess 验证用户访问权限
	ValidateUserAccess(ctx context.Context, userID uint) error

	// GetUserStatus 获取用户状态
	GetUserStatus(ctx context.Context, userID uint) (int, error)

	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, userID uint, loginTime time.Time) error
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}