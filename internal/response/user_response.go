package response

import (
	"time"

	"github.com/LucienVen/golang-auth-service/internal/constants"
	"github.com/LucienVen/golang-auth-service/internal/entity"
)

// UserResponse 用户信息响应
type UserResponse struct {
	ID         string    `json:"id"`         // 用户ID
	Username   string    `json:"username"`   // 用户名
	NickName   string    `json:"nick_name"`  // 昵称
	Email      string    `json:"email"`      // 邮箱
	Phone      string    `json:"phone"`      // 手机号
	Status     int8      `json:"status"`     // 用户状态
	LastLoginAt *int64    `json:"last_login_at"`  // 最后登录时间
	CreatedAt  int64     `json:"created_at"` // 创建时间
	UpdatedAt  int64     `json:"updated_at"` // 更新时间
}

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	User         *UserResponse `json:"user"`         // 用户信息
	AccessToken  string         `json:"access_token"`  // 访问令牌
	RefreshToken string         `json:"refresh_token"` // 刷新令牌
	TokenType    string         `json:"token_type"`    // 令牌类型
	ExpiresIn    int64          `json:"expires_in"`    // 过期时间（秒）
	ExpiresAt    time.Time      `json:"expires_at"`    // 过期时间
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	User         *UserResponse `json:"user"`         // 用户信息
	AccessToken  string         `json:"access_token"`  // 访问令牌
	RefreshToken string         `json:"refresh_token"` // 刷新令牌
	TokenType    string         `json:"token_type"`    // 令牌类型
	ExpiresIn    int64          `json:"expires_in"`    // 过期时间（秒）
	ExpiresAt    time.Time      `json:"expires_at"`    // 过期时间
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	AccessToken  string    `json:"access_token"`  // 访问令牌
	RefreshToken string    `json:"refresh_token"` // 刷新令牌
	TokenType    string    `json:"token_type"`    // 令牌类型
	ExpiresIn    int64     `json:"expires_in"`    // 过期时间（秒）
	ExpiresAt    time.Time `json:"expires_at"`    // 过期时间
}

// UserProfileResponse 用户资料响应
type UserProfileResponse struct {
	ID           string    `json:"id"`           // 用户ID
	Username     string    `json:"username"`     // 用户名
	NickName     string    `json:"nick_name"`    // 昵称
	Email        string    `json:"email"`        // 邮箱
	Phone        string    `json:"phone"`        // 手机号
	Status       int8      `json:"status"`       // 用户状态
	LastLoginAt  *int64    `json:"last_login_at"` // 最后登录时间
	CreatedAt    int64     `json:"created_at"`    // 创建时间
	UpdatedAt    int64     `json:"updated_at"`    // 更新时间
}

// UpdateProfileResponse 更新资料响应
type UpdateProfileResponse struct {
	ID        string    `json:"id"`        // 用户ID
	Username  string    `json:"username"`  // 用户名
	NickName  string    `json:"nick_name"` // 昵称
	Email     string    `json:"email"`     // 邮箱
	Phone     string    `json:"phone"`     // 手机号
	UpdatedAt int64     `json:"updated_at"` // 更新时间
}

// SessionResponse 会话信息响应
type SessionResponse struct {
	SessionID  string    `json:"session_id"`  // 会话ID
	DeviceID   string    `json:"device_id"`   // 设备ID
	UserAgent  string    `json:"user_agent"`  // 用户代理
	IPAddress  string    `json:"ip_address"`  // IP地址
	LoginTime  time.Time `json:"login_time"`  // 登录时间
	LastActive time.Time `json:"last_active"` // 最后活跃时间
	ExpiresAt  time.Time `json:"expires_at"`  // 过期时间
	IsActive   bool      `json:"is_active"`   // 是否活跃
}

// SessionsResponse 用户会话列表响应
type SessionsResponse struct {
	Sessions []*SessionResponse `json:"sessions"` // 会话列表
	Count    int               `json:"count"`    // 会话数量
}

// ChangePasswordResponse 修改密码响应
type ChangePasswordResponse struct {
	UserID    string    `json:"user_id"`    // 用户ID
	UpdatedAt int64     `json:"updated_at"` // 更新时间
	Message   string    `json:"message"`   // 操作结果消息
}

// LogoutResponse 登出响应
type LogoutResponse struct {
	UserID  string `json:"user_id"`  // 用户ID
	Message string `json:"message"`  // 操作结果消息
}

// RevokeSessionResponse 撤销会话响应
type RevokeSessionResponse struct {
	SessionID string `json:"session_id"` // 会话ID
	Message   string `json:"message"`   // 操作结果消息
}

// RevokeAllSessionsResponse 撤销所有会话响应
type RevokeAllSessionsResponse struct {
	Message string `json:"message"` // 操作结果消息
}

// GetStatusText 获取用户状态文本
func GetStatusText(status int8) string {
	switch status {
	case constants.UserStatusInactive:
		return "未激活"
	case constants.UserStatusActive:
		return "正常"
	case constants.UserStatusDisabled:
		return "禁用"
	case constants.UserStatusLogout:
		return "注销"
	case constants.UserStatusDeleted:
		return "已删除"
	default:
		return "未知"
	}
}

// ToUserResponse 将用户实体转换为用户响应
func ToUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		NickName:   user.NickName,
		Email:      user.Email,
		Phone:      user.Phone,
		Status:     user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:  user.CreateTime,
		UpdatedAt:  user.UpdateTime,
	}
}

// ToUserProfileResponse 将用户实体转换为用户资料响应
func ToUserProfileResponse(user *entity.User) *UserProfileResponse {
	return &UserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		NickName:    user.NickName,
		Email:       user.Email,
		Phone:       user.Phone,
		Status:      user.Status,
		LastLoginAt:  user.LastLoginAt,
		CreatedAt:    user.CreateTime,
		UpdatedAt:    user.UpdateTime,
	}
}

// ToSessionResponse 将会话信息转换为响应
func ToSessionResponse(session interface{}) *SessionResponse {
	// 类型断言，支持不同的会话信息类型
	switch s := session.(type) {
	case map[string]interface{}:
		// 处理 map 类型的会话信息
		return &SessionResponse{
			SessionID:  getStringFromMap(s, "session_id"),
			DeviceID:   getStringFromMap(s, "device_id"),
			UserAgent:  getStringFromMap(s, "user_agent"),
			IPAddress:  getStringFromMap(s, "ip_address"),
			LoginTime:  getTimeFromMap(s, "login_time"),
			LastActive: getTimeFromMap(s, "last_active"),
			ExpiresAt:  getTimeFromMap(s, "expires_at"),
			IsActive:   getBoolFromMap(s, "is_active"),
		}
	default:
		// 处理其他类型，返回空的响应
		return &SessionResponse{}
	}
}

// 辅助函数：从 map 中安全获取字符串值
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// 辅助函数：从 map 中安全获取时间值
func getTimeFromMap(m map[string]interface{}, key string) time.Time {
	if val, ok := m[key]; ok {
		if t, ok := val.(time.Time); ok {
			return t
		}
	}
	return time.Time{}
}

// 辅助函数：从 map 中安全获取布尔值
func getBoolFromMap(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}