package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/entity"
	"github.com/LucienVen/golang-auth-service/internal/validator"
)

// GetUserStatus 获取用户状态实现
func (s *UserServiceImpl) GetUserStatus(ctx context.Context, userID uint) (int, error) {
	// 获取用户状态
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("获取用户状态失败: %w", err)
	}

	if user == nil {
		return 0, validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 检查是否被逻辑删除
	if user.IsDelete > 0 {
		return int(entity.UserStatusDeleted), nil
	}

	return int(user.Status), nil
}

// UpdateLastLogin 更新最后登录时间实现
func (s *UserServiceImpl) UpdateLastLogin(ctx context.Context, userID uint, loginTime time.Time) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user == nil {
		return validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 检查用户状态
	if user.IsDelete > 0 {
		return validator.NewValidationError("user_id", "用户已被删除", "user_deleted")
	}

	if user.Status == entity.UserStatusDisabled { // 禁用状态
		return validator.NewValidationError("user_id", "用户已被禁用", "user_disabled")
	}

	// 更新最后登录时间
	loginUnix := loginTime.Unix()
	user.LastLoginAt = &loginUnix
	user.Updater = "login_service"

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("更新最后登录时间失败: %w", err)
	}

	return nil
}

// DeactivateUser 停用用户（软删除）实现
func (s *UserServiceImpl) DeactivateUser(ctx context.Context, userID uint) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user == nil {
		return validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 检查用户是否已经被停用
	if user.IsDelete > 0 {
		return validator.NewValidationError("user_id", "用户已被删除", "user_deleted")
	}

	// 执行逻辑删除
	user.IsDelete = 1
	user.Status = entity.UserStatusDeleted
	user.Updater = "deactivate_service"

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("停用用户失败: %w", err)
	}

	return nil
}

// ActivateUser 激活用户（管理员功能）实现
func (s *UserServiceImpl) ActivateUser(ctx context.Context, userID uint) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user == nil {
		return validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 激活用户
	user.Status = entity.UserStatusActive
	user.IsDelete = 0 // 清除逻辑删除标记
	user.Updater = "activate_service"

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("激活用户失败: %w", err)
	}

	return nil
}

// DisableUser 禁用用户（管理员功能）实现
func (s *UserServiceImpl) DisableUser(ctx context.Context, userID uint) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user == nil {
		return validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 禁用用户
	user.Status = entity.UserStatusDisabled
	user.Updater = "disable_service"

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("禁用用户失败: %w", err)
	}

	return nil
}

// IsUserActive 检查用户是否处于活跃状态
func (s *UserServiceImpl) IsUserActive(ctx context.Context, userID uint) (bool, error) {
	status, err := s.GetUserStatus(ctx, userID)
	if err != nil {
		return false, err
	}

	return status == int(entity.UserStatusActive), nil
}

// ValidateUserAccess 验证用户访问权限
func (s *UserServiceImpl) ValidateUserAccess(ctx context.Context, userID uint) error {
	// 检查用户是否存在且状态正常
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("验证用户访问权限失败: %w", err)
	}

	if user == nil {
		return validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 检查用户状态
	switch user.Status {
	case entity.UserStatusInactive:
		return validator.NewValidationError("user_id", "用户未激活", "user_not_activated")
	case entity.UserStatusDisabled:
		return validator.NewValidationError("user_id", "用户已被禁用", "user_disabled")
	case entity.UserStatusDeleted:
		return validator.NewValidationError("user_id", "用户已注销", "user_deleted")
	}

	// 检查逻辑删除状态
	if user.IsDelete > 0 {
		return validator.NewValidationError("user_id", "用户已被删除", "user_deleted")
	}

	return nil
}