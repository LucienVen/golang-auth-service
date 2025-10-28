package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/entity"
	"github.com/LucienVen/golang-auth-service/internal/request"
	"github.com/LucienVen/golang-auth-service/internal/response"
	"github.com/LucienVen/golang-auth-service/internal/validator"
)

// GetUserProfile 获取用户信息实现
func (s *UserServiceImpl) GetUserProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error) {
	// 1. 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user == nil {
		return nil, validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 2. 检查用户状态（逻辑删除的用户不应该返回信息）
	if user.IsDelete > 0 {
		return nil, validator.NewValidationError("user_id", "用户已被删除", "user_deleted")
	}

	// 3. 构造响应
	userIDUint, _ := strconv.ParseUint(user.ID, 10, 32)
	lastLogin := int64(0)
	if user.LastLoginAt != nil {
		lastLogin = *user.LastLoginAt
	}

	return &response.UserProfileResponse{
		UserID:    uint(userIDUint),
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		NickName:  user.NickName,
		Status:    int(user.Status),
		LastLogin: lastLogin,
		CreatedAt: time.Unix(user.CreateTime, 0),
		UpdatedAt: time.Unix(user.UpdateTime, 0),
	}, nil
}

// GetUserProfileByAccount 通过账号获取用户信息（内部使用）
func (s *UserServiceImpl) GetUserProfileByAccount(ctx context.Context, account string) (*response.UserProfileResponse, error) {
	if account == "" {
		return nil, validator.NewValidationError("account", "账号不能为空", "required")
	}

	// 1. 通过统一账号查询获取用户
	user, err := s.userRepo.GetByAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("通过账号获取用户信息失败: %w", err)
	}

	if user == nil {
		return nil, validator.NewValidationError("account", "用户不存在", "user_not_found")
	}

	// 2. 检查用户状态
	if user.IsDelete > 0 {
		return nil, validator.NewValidationError("account", "用户已被删除", "user_deleted")
	}

	// 3. 构造响应
	userIDUint, _ := strconv.ParseUint(user.ID, 10, 32)
	lastLogin := int64(0)
	if user.LastLoginAt != nil {
		lastLogin = *user.LastLoginAt
	}

	return &response.UserProfileResponse{
		UserID:    uint(userIDUint),
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		NickName:  user.NickName,
		Status:    int(user.Status),
		LastLogin: lastLogin,
		CreatedAt: time.Unix(user.CreateTime, 0),
		UpdatedAt: time.Unix(user.UpdateTime, 0),
	}, nil
}

// UpdateProfile 更新用户资料实现
func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uint, req *request.UpdateProfileRequest) (*response.UpdateProfileResponse, error) {
	// 1. 基础验证
	if validationResult := validator.ValidateAndFormat(req); !validationResult.Valid {
		return nil, fmt.Errorf("输入验证失败: %v", validationResult.ToMap())
	}

	// 2. 业务规则验证 - 至少提供一个更新字段
	if req.NickName == "" && req.Email == "" && req.Phone == "" {
		return nil, validator.NewValidationError("profile", "至少需要提供一个要更新的字段", "required_field")
	}

	// 3. 获取当前用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user == nil {
		return nil, validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 4. 检查用户状态
	if user.IsDelete > 0 {
		return nil, validator.NewValidationError("user_id", "用户已被删除", "user_deleted")
	}

	if user.Status == entity.UserStatusDisabled { // 禁用状态
		return nil, validator.NewValidationError("user_id", "用户已被禁用", "user_disabled")
	}

	// 5. 验证更新账号的可用性
	if err := s.ValidateUpdateAccountAvailable(ctx, userID, req); err != nil {
		return nil, err
	}

	// 6. 更新字段
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	user.Updater = "update_profile_service"

	// 7. 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("更新用户信息失败: %w", err)
	}

	// 8. 构造响应
	userIDUint, _ := strconv.ParseUint(user.ID, 10, 32)
	return &response.UpdateProfileResponse{
		UserID:    uint(userIDUint),
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		NickName:  user.NickName,
		Status:    int(user.Status),
		UpdatedAt: time.Unix(user.UpdateTime, 0),
	}, nil
}