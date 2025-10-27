package entity

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/utils"
	"gorm.io/gorm"
)

// UserStatusInactive 表示用户未激活。
// UserStatusActive 表示用户正常。
// UserStatusDisabled 表示用户被禁用。
// UserStatusLogout 表示用户已注销。
// UserStatusDeleted 表示用户已删除。
const (
	UserStatusInactive int8 = 0 // 未激活
	UserStatusActive   int8 = 1 // 正常
	UserStatusDisabled int8 = 2 // 禁用
	UserStatusLogout   int8 = 3 // 注销
	UserStatusDeleted  int8 = 9 // 已删除
)

// User 用户实体
type User struct {
	ID           string `json:"id" gorm:"column:id;primaryKey;autoIncrement"`                    // 用户ID
	Username     string `json:"username" gorm:"column:username;uniqueIndex;size:50"`            // 用户名
	NickName     string `json:"nick_name" gorm:"column:nick_name;size:100"`                     // 昵称
	PasswordHash string `json:"-" gorm:"column:password_hash;size:255;not null"`                // 密码哈希（已加密）
	Phone        string `json:"phone" gorm:"column:phone;uniqueIndex;size:20"`                  // 手机号
	Email        string `json:"email" gorm:"column:email;uniqueIndex;size:100"`                 // 邮箱地址
	Status       int8   `json:"status" gorm:"column:status;default:0;index"`                    // 用户状态：0-未激活，1-正常，2-禁用，3-注销，9-已删除
	LastLoginAt  *int64 `json:"last_login_at" gorm:"column:last_login_at"`                      // 最后登录时间
	BaseModel
}

// TableName 指定表名
func (u *User) TableName() string {
	return UserTable
}

// NewUser 创建新用户
func NewUser(username, nickName, password, phone, email string) (*User, error) {
	user := &User{
		BaseModel: GenBaseModel("system", "system"),
		Username:  username,
		NickName:  nickName,
		Phone:     phone,
		Email:     email,
		Status:    UserStatusInactive,
	}

	// 验证并设置密码
	if err := user.ValidatePassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

// NewUserWithHash 创建新用户（使用已有密码哈希）
func NewUserWithHash(username, nickName, passwordHash, phone, email string) *User {
	return &User{
		BaseModel:    GenBaseModel("system", "system"),
		Username:     username,
		NickName:     nickName,
		PasswordHash: passwordHash,
		Phone:        phone,
		Email:        email,
		Status:       UserStatusInactive,
	}
}

// 用户状态事件常量
const (
	UserEventActivate = "activate" // 激活
	UserEventDisable  = "disable"  // 封号
	UserEventLogout   = "logout"   // 注销
	UserEventDelete   = "delete"   // 删除
)

// 用户状态转移表
var userStatusTransitions = map[int8]map[string]int8{
	UserStatusInactive: {
		UserEventActivate: UserStatusActive, // 激活
	},
	UserStatusActive: {
		UserEventDisable: UserStatusDisabled, // 封号
		UserEventLogout:  UserStatusLogout,   // 注销
		UserEventDelete:  UserStatusDeleted,  // 删除
	},
	UserStatusDisabled: {
		UserEventActivate: UserStatusActive,  // 解封
		UserEventDelete:   UserStatusDeleted, // 删除
	},
	UserStatusLogout: {
		UserEventActivate: UserStatusActive,  // 恢复
		UserEventDelete:   UserStatusDeleted, // 删除
	},
}

// TransitionUserStatus 用户状态转移
// event: 事件名，建议使用 UserEventActivate/UserEventDisable/UserEventLogout/UserEventDelete
// 返回 true 表示转移成功，false 表示非法转移
func (u *User) TransitionUserStatus(event string) bool {
	if next, ok := userStatusTransitions[u.Status][event]; ok {
		u.Status = next
		return true
	}
	return false
}

// ========== 状态检查方法 ==========

// IsActive 检查用户是否处于正常状态
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsInactive 检查用户是否未激活
func (u *User) IsInactive() bool {
	return u.Status == UserStatusInactive
}

// IsDisabled 检查用户是否被禁用
func (u *User) IsDisabled() bool {
	return u.Status == UserStatusDisabled
}

// IsLogout 检查用户是否已注销
func (u *User) IsLogout() bool {
	return u.Status == UserStatusLogout
}

// IsDeleted 检查用户是否已删除
func (u *User) IsDeleted() bool {
	return u.Status == UserStatusDeleted
}

// CanLogin 检查用户是否可以登录
func (u *User) CanLogin() bool {
	return u.IsActive()
}

// ========== 验证方法 ==========

// ValidateUsername 验证用户名格式
func (u *User) ValidateUsername() error {
	if u.Username == "" {
		return errors.New("用户名不能为空")
	}

	if len(u.Username) < 3 || len(u.Username) > 50 {
		return errors.New("用户名长度必须在3-50个字符之间")
	}

	// 用户名只能包含字母、数字、下划线，且不能以下划线开头或结尾
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9_]*[a-zA-Z0-9]$`, u.Username)
	if !matched {
		return errors.New("用户名只能包含字母、数字、下划线，且不能以下划线开头或结尾")
	}

	return nil
}

// ValidateEmail 验证邮箱格式
func (u *User) ValidateEmail() error {
	if u.Email == "" {
		return errors.New("邮箱不能为空")
	}

	// 基本邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("邮箱格式不正确")
	}

	// 邮箱长度限制
	if len(u.Email) > 100 {
		return errors.New("邮箱长度不能超过100个字符")
	}

	return nil
}

// ValidatePhone 验证手机号格式（支持中国大陆）
func (u *User) ValidatePhone() error {
	if u.Phone == "" {
		return errors.New("手机号不能为空")
	}

	// 中国大陆手机号格式验证：1开头，第二位3-9，共11位数字
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(u.Phone) {
		return errors.New("手机号格式不正确，请输入有效的中国大陆手机号")
	}

	return nil
}

// ValidatePassword 验证密码强度并设置哈希
func (u *User) ValidatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}

	// 使用密码工具验证强度
	if err := utils.CheckPasswordStrength(password); err != nil {
		return err
	}

	// 生成密码哈希
	hash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	u.PasswordHash = hash
	return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) error {
	if u.PasswordHash == "" {
		return errors.New("用户密码未设置")
	}

	return utils.VerifyPassword(password, u.PasswordHash)
}

// ========== 便利方法 ==========

// GetDisplayUsername 获取显示用户名（优先级：昵称 > 用户名）
func (u *User) GetDisplayUsername() string {
	if u.NickName != "" {
		return u.NickName
	}
	return u.Username
}

// GetAccountIdentifier 获取账户标识符（用于登录）
func (u *User) GetAccountIdentifier() string {
	if u.Email != "" {
		return u.Email
	}
	if u.Phone != "" {
		return u.Phone
	}
	return u.Username
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := time.Now().Unix()
	u.LastLoginAt = &now
	u.UpdateTime = now
}

// Sanitize 清理和标准化用户数据
func (u *User) Sanitize() {
	u.Username = strings.TrimSpace(strings.ToLower(u.Username))
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.Phone = strings.TrimSpace(u.Phone)
	u.NickName = strings.TrimSpace(u.NickName)
}

// ToSafeJSON 转换为安全的JSON格式（排除敏感信息）
func (u *User) ToSafeJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":             u.ID,
		"username":       u.Username,
		"nick_name":      u.NickName,
		"email":          u.Email,
		"phone":          u.Phone,
		"status":         u.Status,
		"last_login_at":  u.LastLoginAt,
		"create_time":    u.CreateTime,
		"update_time":    u.UpdateTime,
	}
}

// ========== GORM 钩子 ==========

// BeforeCreate GORM创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 清理数据
	u.Sanitize()

	// 验证必填字段
	if u.Username == "" {
		return errors.New("用户名不能为空")
	}

	if u.PasswordHash == "" {
		return errors.New("密码不能为空")
	}

	// 如果邮箱不为空，验证格式
	if u.Email != "" {
		if err := u.ValidateEmail(); err != nil {
			return err
		}
	}

	// 如果手机号不为空，验证格式
	if u.Phone != "" {
		if err := u.ValidatePhone(); err != nil {
			return err
		}
	}

	// 验证用户名格式
	if err := u.ValidateUsername(); err != nil {
		return err
	}

	// 设置默认状态
	if u.Status == 0 {
		u.Status = UserStatusInactive
	}

	return nil
}

// BeforeUpdate GORM更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 清理数据
	u.Sanitize()

	// 验证格式（如果字段有更新）
	if u.Username != "" {
		if err := u.ValidateUsername(); err != nil {
			return err
		}
	}

	if u.Email != "" {
		if err := u.ValidateEmail(); err != nil {
			return err
		}
	}

	if u.Phone != "" {
		if err := u.ValidatePhone(); err != nil {
			return err
		}
	}

	return nil
}
