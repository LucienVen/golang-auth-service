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
	DeviceID string `json:"device_id" binding:"omitempty"`                    // 设备ID
	UserAgent string `json:"user_agent" binding:"omitempty"`                    // 用户代理
	IPAddress string `json:"ip_address" binding:"omitempty"`                    // IP地址
}

// LoginRequest 用户登录请求（统一账号识别）
type LoginRequest struct {
	Account   string `json:"account" binding:"required"`   // 账号（邮箱/手机号/用户名）
	Password  string `json:"password" binding:"required"`                 // 密码
	DeviceID  string `json:"device_id" binding:"omitempty"`              // 设备ID
	UserAgent string `json:"user_agent" binding:"omitempty"`              // 用户代理
	IPAddress string `json:"ip_address" binding:"omitempty"`              // IP地址
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

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`  // 刷新令牌
}

// RevokeSessionRequest 撤销会话请求
type RevokeSessionRequest struct {
	SessionID string `json:"session_id" binding:"required"`  // 会话ID
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`  // 令牌
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

// GetAccount 获取注册请求的主要账号（用于账号唯一性检查）
func (r *RegisterRequest) GetAccount() string {
	if r.Email != "" {
		return r.Email
	}
	if r.Phone != "" {
		return r.Phone
	}
	return r.Username
}

// Validate 注册请求整体验证
func (r *RegisterRequest) Validate() error {
	if err := r.ValidateAccountType(); err != nil {
		return err
	}
	return nil
}

// Validate 登录请求整体验证
func (r *LoginRequest) Validate() error {
	return r.ValidateAccount()
}