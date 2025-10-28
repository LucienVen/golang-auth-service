package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/entity"
	"github.com/LucienVen/golang-auth-service/internal/request"
	"github.com/LucienVen/golang-auth-service/internal/response"
	"github.com/LucienVen/golang-auth-service/internal/utils"
	"github.com/LucienVen/golang-auth-service/internal/validator"
)

// ChangePassword 修改密码实现
func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) (*response.ChangePasswordResponse, error) {
	// 1. 基础验证
	if validationResult := validator.ValidateAndFormat(req); !validationResult.Valid {
		return nil, fmt.Errorf("输入验证失败: %v", validationResult.ToMap())
	}

	// 2. 获取当前用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	if user == nil {
		return nil, validator.NewValidationError("user_id", "用户不存在", "user_not_found")
	}

	// 3. 检查用户状态
	if user.IsDelete > 0 {
		return nil, validator.NewValidationError("user_id", "用户已被删除", "user_deleted")
	}

	if user.Status == entity.UserStatusDisabled { // 禁用状态
		return nil, validator.NewValidationError("user_id", "用户已被禁用", "user_disabled")
	}

	// 4. 验证旧密码
	if err := utils.VerifyPassword(req.OldPassword, user.PasswordHash); err != nil {
		return nil, validator.NewValidationError("old_password", "旧密码不正确", "password_incorrect")
	}

	// 5. 新密码强度检查已在validator中完成
	// 6. 新旧密码不能相同检查已在validator中完成

	// 7. 生成新密码哈希
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("新密码哈希生成失败: %w", err)
	}

	// 8. 更新密码
	user.PasswordHash = newPasswordHash
	user.Updater = "change_password_service"

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("更新密码失败: %w", err)
	}

	// 9. 构造响应
	userIDUint, _ := strconv.ParseUint(user.ID, 10, 32)
	return &response.ChangePasswordResponse{
		UserID:    uint(userIDUint),
		UpdatedAt: time.Unix(user.UpdateTime, 0),
		Message:   "密码修改成功",
	}, nil
}

// ValidateAccountExists 检查账号是否存在实现
func (s *UserServiceImpl) ValidateAccountExists(ctx context.Context, req *request.RegisterRequest) error {
	var errors []validator.ValidationError

	// 检查用户名是否已存在
	if req.Username != "" {
		if user, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil && user != nil {
			errors = append(errors, validator.NewValidationError("username", "用户名已存在", "account_exists"))
		}
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if user, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil && user != nil {
			errors = append(errors, validator.NewValidationError("email", "邮箱已存在", "account_exists"))
		}
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		if user, err := s.userRepo.GetByPhone(ctx, req.Phone); err == nil && user != nil {
			errors = append(errors, validator.NewValidationError("phone", "手机号已存在", "account_exists"))
		}
	}

	// 如果有重复账号，返回错误
	if len(errors) > 0 {
		validationResult := validator.ValidateBusiness(errors...)
		if !validationResult.Valid {
			return fmt.Errorf("账号验证失败: %v", validationResult.ToMap())
		}
	}

	return nil
}

// ValidateUpdateAccountAvailable 验证更新账号的可用性实现
func (s *UserServiceImpl) ValidateUpdateAccountAvailable(ctx context.Context, userID uint, req *request.UpdateProfileRequest) error {
	var errors []validator.ValidationError

	// 检查邮箱是否被其他用户使用
	if req.Email != "" {
		if user, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil && user != nil {
			userIDUint, _ := strconv.ParseUint(user.ID, 10, 32)
			if uint(userIDUint) != userID {
				errors = append(errors, validator.NewValidationError("email", "邮箱已被其他用户使用", "account_exists"))
			}
		}
	}

	// 检查手机号是否被其他用户使用
	if req.Phone != "" {
		if user, err := s.userRepo.GetByPhone(ctx, req.Phone); err == nil && user != nil {
			userIDUint, _ := strconv.ParseUint(user.ID, 10, 32)
			if uint(userIDUint) != userID {
				errors = append(errors, validator.NewValidationError("phone", "手机号已被其他用户使用", "account_exists"))
			}
		}
	}

	// 如果有冲突账号，返回错误
	if len(errors) > 0 {
		validationResult := validator.ValidateBusiness(errors...)
		if !validationResult.Valid {
			return fmt.Errorf("账号可用性验证失败: %v", validationResult.ToMap())
		}
	}

	return nil
}