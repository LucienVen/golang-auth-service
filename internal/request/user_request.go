package request

import (
	"github.com/LucienVen/golang-auth-service/internal/validator"
)

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"omitempty,username_valid"`        // 用户名（可选）
	Email    string `json:"email" binding:"omitempty,email"`                    // 邮箱（可选，与phone至少一个）
	Phone    string `json:"phone" binding:"omitempty,phone_valid"`              // 手机号（可选，与email至少一个）
	NickName string `json:"nick_name" binding:"omitempty,max=100"`              // 昵称（可选）
	Password string `json:"password" binding:"required,strong_password"`        // 密码
}

// LoginRequest 用户登录请求（统一账号识别）
type LoginRequest struct {
	Account  string `json:"account" binding:"required,account_type"`   // 账号（邮箱/手机号/用户名）
	Password string `json:"password" binding:"required"`                 // 密码
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	NickName string `json:"nick_name" binding:"omitempty,max=100"`      // 昵称
	Email    string `json:"email" binding:"omitempty,email"`            // 邮箱
	Phone    string `json:"phone" binding:"omitempty,phone_valid"`      // 手机号
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`             // 旧密码
	NewPassword string `json:"new_password" binding:"required,strong_password"` // 新密码
}

// GetUserProfileRequest 获取用户信息请求（预留）
type GetUserProfileRequest struct {
	UserID uint `json:"user_id" binding:"required"`  // 用户ID
}

// 用户请求验证方法

// ValidateAccountType 判断账号类型
func (r *RegisterRequest) ValidateAccountType() error {
	if r.Email == "" && r.Phone == "" && r.Username == "" {
		return validator.NewValidationError("username", "用户名、邮箱、手机号至少需要提供一个", "required_account")
	}
	return nil
}

// ValidateAccount 判断登录账号是否有效
func (r *LoginRequest) ValidateAccount() error {
	if len(r.Account) == 0 {
		return validator.NewValidationError("account", "账号不能为空", "required")
	}
	if len(r.Account) > 100 {
		return validator.NewValidationError("account", "账号长度不能超过100个字符", "max_length")
	}
	return nil
}

// ValidateUpdateProfile 验证更新资料请求
func (r *UpdateProfileRequest) ValidateUpdateProfile() error {
	if r.Email == "" && r.Phone == "" && r.NickName == "" {
		return validator.NewValidationError("profile", "至少需要提供一个要更新的字段", "required_field")
	}
	return nil
}

// ValidateChangePassword 验证修改密码请求
func (r *ChangePasswordRequest) ValidateChangePassword() error {
	if r.OldPassword == r.NewPassword {
		return validator.NewValidationError("new_password", "新密码不能与旧密码相同", "same_as_old")
	}
	return nil
}