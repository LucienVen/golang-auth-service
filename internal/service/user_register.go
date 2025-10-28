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

// Register 用户注册实现
func (s *UserServiceImpl) Register(ctx context.Context, req *request.RegisterRequest) (*response.RegisterResponse, error) {
	// 1. 基础验证
	if validationResult := validator.ValidateAndFormat(req); !validationResult.Valid {
		return nil, fmt.Errorf("输入验证失败: %v", validationResult.ToMap())
	}

	// 2. 业务规则验证 - 至少提供一个账号
	if req.Username == "" && req.Email == "" && req.Phone == "" {
		return nil, validator.NewValidationError("account", "用户名、邮箱、手机号至少需要提供一个", "required_account")
	}

	// 3. 检查账号是否已存在
	if err := s.ValidateAccountExists(ctx, req); err != nil {
		return nil, err
	}

	// 4. 密码强度检查已在validator中完成

	// 5. 创建用户实体
	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		NickName:     req.NickName,
		PasswordHash: "", // 将在BeforeCreate钩子中设置
		Status:       entity.UserStatusInactive, // 默认未激活状态
	}

	// 6. 设置密码哈希（在BeforeCreate钩子之前手动设置）
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希生成失败: %w", err)
	}
	user.PasswordHash = passwordHash

	// 7. 保存到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 8. 构造响应
	userIDUint, _ := strconv.ParseUint(user.ID, 10, 32)
	return &response.RegisterResponse{
		UserID:    uint(userIDUint),
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		NickName:  user.NickName,
		Status:    int(user.Status),
		CreatedAt: time.Unix(user.CreateTime, 0),
	}, nil
}