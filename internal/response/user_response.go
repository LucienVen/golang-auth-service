package response

import (
	"time"
)

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	UserID    uint      `json:"user_id"`    // 用户ID
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱
	Phone     string    `json:"phone"`      // 手机号
	NickName  string    `json:"nick_name"`  // 昵称
	Status    int       `json:"status"`     // 用户状态
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	UserID    uint      `json:"user_id"`    // 用户ID
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱
	Phone     string    `json:"phone"`      // 手机号
	NickName  string    `json:"nick_name"`  // 昵称
	Status    int       `json:"status"`     // 用户状态
	Token     string    `json:"token"`      // JWT令牌
	ExpiresAt time.Time `json:"expires_at"` // 令牌过期时间
	LastLogin time.Time `json:"last_login"` // 最后登录时间
}

// UserProfileResponse 用户资料响应
type UserProfileResponse struct {
	UserID     uint      `json:"user_id"`     // 用户ID
	Username   string    `json:"username"`    // 用户名
	Email      string    `json:"email"`       // 邮箱
	Phone      string    `json:"phone"`       // 手机号
	NickName   string    `json:"nick_name"`   // 昵称
	Status     int       `json:"status"`      // 用户状态
	LastLogin  int64     `json:"last_login"`  // 最后登录时间
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}

// UpdateProfileResponse 更新资料响应
type UpdateProfileResponse struct {
	UserID    uint      `json:"user_id"`    // 用户ID
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱
	Phone     string    `json:"phone"`      // 手机号
	NickName  string    `json:"nick_name"`  // 昵称
	Status    int       `json:"status"`     // 用户状态
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// ChangePasswordResponse 修改密码响应
type ChangePasswordResponse struct {
	UserID    uint      `json:"user_id"`    // 用户ID
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
	Message   string    `json:"message"`   // 操作结果消息
}

// LogoutResponse 登出响应
type LogoutResponse struct {
	UserID  uint   `json:"user_id"`  // 用户ID
	Message string `json:"message"`  // 操作结果消息
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	UserID    uint      `json:"user_id"`    // 用户ID
	Token     string    `json:"token"`      // 新JWT令牌
	ExpiresAt time.Time `json:"expires_at"` // 令牌过期时间
}

// UserStatusResponse 用户状态响应
type UserStatusResponse struct {
	UserID      uint      `json:"user_id"`      // 用户ID
	Username    string    `json:"username"`     // 用户名
	Status      int       `json:"status"`      // 用户状态
	StatusText  string    `json:"status_text"`  // 状态描述
	LastLogin   int64     `json:"last_login"`   // 最后登录时间
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64                        `json:"total"` // 总数
	Items []*UserProfileResponse        `json:"items"` // 用户列表
}

// CreateUserResponse 创建用户响应（管理员使用）
type CreateUserResponse struct {
	UserID    uint      `json:"user_id"`    // 用户ID
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱
	Phone     string    `json:"phone"`      // 手机号
	NickName  string    `json:"nick_name"`  // 昵称
	Status    int       `json:"status"`     // 用户状态
	CreatedAt time.Time `json:"created_at"` // 创建时间
	Message   string    `json:"message"`   // 操作结果消息
}

// UpdateUserResponse 更新用户响应（管理员使用）
type UpdateUserResponse struct {
	UserID    uint      `json:"user_id"`    // 用户ID
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱
	Phone     string    `json:"phone"`      // 手机号
	NickName  string    `json:"nick_name"`  // 昵称
	Status    int       `json:"status"`     // 用户状态
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
	Message   string    `json:"message"`   // 操作结果消息
}

// DeleteUserResponse 删除用户响应（管理员使用）
type DeleteUserResponse struct {
	UserID  uint   `json:"user_id"`  // 用户ID
	Message string `json:"message"`  // 操作结果消息
}

// EnableUserResponse 启用用户响应（管理员使用）
type EnableUserResponse struct {
	UserID  uint   `json:"user_id"`  // 用户ID
	Message string `json:"message"`  // 操作结果消息
}

// DisableUserResponse 禁用用户响应（管理员使用）
type DisableUserResponse struct {
	UserID  uint   `json:"user_id"`  // 用户ID
	Message string `json:"message"`  // 操作结果消息
}

// UserStatsResponse 用户统计响应
type UserStatsResponse struct {
	TotalUsers   int64 `json:"total_users"`   // 总用户数
	ActiveUsers  int64 `json:"active_users"`  // 活跃用户数
	InactiveUsers int64 `json:"inactive_users"` // 未激活用户数
	DisabledUsers int64 `json:"disabled_users"` // 禁用用户数
	TodayUsers   int64 `json:"today_users"`   // 今日注册用户数
	WeekUsers    int64 `json:"week_users"`    // 本周注册用户数
	MonthUsers   int64 `json:"month_users"`   // 本月注册用户数
}

// 用户状态常量
const (
	UserStatusUnactivated = 0 // 未激活
	UserStatusActive      = 1 // 正常
	UserStatusDisabled   = 2 // 禁用
	UserStatusDeleted    = 3 // 注销
	UserStatusDeletedLogically = 9 // 已删除（软删除）
)

// GetStatusText 获取用户状态文本
func GetStatusText(status int) string {
	switch status {
	case UserStatusUnactivated:
		return "未激活"
	case UserStatusActive:
		return "正常"
	case UserStatusDisabled:
		return "禁用"
	case UserStatusDeleted:
		return "注销"
	case UserStatusDeletedLogically:
		return "已删除"
	default:
		return "未知"
	}
}