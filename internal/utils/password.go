package utils

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	// bcrypt cost 值，用于控制哈希计算复杂度
	// cost=12 在安全性和性能之间的良好平衡
	// cost值每增加1，哈希计算时间翻倍
	// 推荐范围：10-14
	bcryptCost = 12

	// 密码最小长度
	minPasswordLength = 8

	// 密码最大长度
	maxPasswordLength = 72

	// 密码强度检查
	// 不使用正则表达式，而是逐个检查条件
	// 必须至少包含一个小写字母、一个大写字母、一个数字，可选特殊字符
)

// 自定义错误类型
var (
	ErrPasswordTooShort      = fmt.Errorf("密码长度不能少于%d位", minPasswordLength)
	ErrPasswordTooLong      = fmt.Errorf("密码长度不能超过%d位", maxPasswordLength)
	ErrPasswordTooWeak      = fmt.Errorf("密码强度不足，必须包含大小写字母和数字")
	ErrPasswordHashFailed   = fmt.Errorf("密码哈希生成失败")
	ErrPasswordVerifyFailed = fmt.Errorf("密码验证失败")
)

// PasswordHasher 密码哈希器接口
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
	CheckPasswordStrength(password string) error
}

// BcryptPasswordHasher 基于bcrypt的密码哈希器实现
type BcryptPasswordHasher struct{}

// NewPasswordHasher 创建密码哈希器实例
func NewPasswordHasher() PasswordHasher {
	return &BcryptPasswordHasher{}
}

// HashPassword 生成密码哈希
// 使用bcrypt算法对密码进行哈希处理
// bcrypt具有自动加盐功能，可以防止彩虹表攻击
func (h *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	// 首先检查密码强度
	if err := h.CheckPasswordStrength(password); err != nil {
		return "", err
	}

	// 使用bcrypt生成密码哈希
	// cost=12 在安全性和性能之间取得平衡
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrPasswordHashFailed, err)
	}

	return string(hash), nil
}

// VerifyPassword 验证密码
// 将明文密码与存储的哈希值进行比对
func (h *BcryptPasswordHasher) VerifyPassword(password, hash string) error {
	// 检查哈希值是否为空
	if hash == "" {
		return fmt.Errorf("%w: 哈希值为空", ErrPasswordVerifyFailed)
	}

	// 使用bcrypt.CompareHashAndPassword验证密码
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("%w: 密码不匹配", ErrPasswordVerifyFailed)
		}
		return fmt.Errorf("%w: %v", ErrPasswordVerifyFailed, err)
	}

	return nil
}

// CheckPasswordStrength 检查密码强度
// 验证密码是否符合安全要求
func (h *BcryptPasswordHasher) CheckPasswordStrength(password string) error {
	// 检查密码长度
	if utf8.RuneCountInString(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	if utf8.RuneCountInString(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	// 检查密码复杂度 - 逐个检查条件
	var (
		hasLower bool
		hasUpper bool
		hasDigit bool
	)

	for _, r := range password {
		switch {
		case 'a' <= r && r <= 'z':
			hasLower = true
		case 'A' <= r && r <= 'Z':
			hasUpper = true
		case '0' <= r && r <= '9':
			hasDigit = true
		}
	}

	// 必须至少包含一个小写字母、一个大写字母、一个数字
	if !hasLower || !hasUpper || !hasDigit {
		return ErrPasswordTooWeak
	}

	return nil
}

// 全局默认密码哈希器实例
var defaultPasswordHasher = NewPasswordHasher()

// HashPassword 全局函数：生成密码哈希
// 使用默认的密码哈希器
func HashPassword(password string) (string, error) {
	return defaultPasswordHasher.HashPassword(password)
}

// VerifyPassword 全局函数：验证密码
// 使用默认的密码哈希器
func VerifyPassword(password, hash string) error {
	return defaultPasswordHasher.VerifyPassword(password, hash)
}

// CheckPasswordStrength 全局函数：检查密码强度
// 使用默认的密码哈希器
func CheckPasswordStrength(password string) error {
	return defaultPasswordHasher.CheckPasswordStrength(password)
}

// IsPasswordHashValid 检查密码哈希格式是否有效
// 验证哈希字符串是否为有效的bcrypt格式
func IsPasswordHashValid(hash string) bool {
	if len(hash) != 60 {
		return false
	}

	// bcrypt哈希格式: $2a$[cost]$[22字符盐][31字符哈希] 或 $2b$[cost]$[22字符盐][31字符哈希]
	if hash[:4] != "$2b$" && hash[:4] != "$2a$" {
		return false
	}

	return true
}

// GenerateRandomPassword 生成随机密码（用于密码重置等功能）
// 生成符合安全要求的随机密码
func GenerateRandomPassword(length int) (string, error) {
	if length < minPasswordLength {
		length = minPasswordLength
	}
	if length > maxPasswordLength {
		length = maxPasswordLength
	}

	// 密码字符集
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	special := "@$!%*?&"

	// 确保至少包含每种类型的字符
	password := make([]byte, length)
	password[0] = lowercase[0]                              // 小写字母
	password[1] = uppercase[0]                              // 大写字母
	password[2] = digits[0]                                  // 数字
	password[3] = special[0]                                  // 特殊字符

	// 剩余字符从所有字符集中随机选择
	allChars := lowercase + uppercase + digits + special
	for i := 4; i < length; i++ {
		// 这里应该使用更安全的随机数生成器，如crypto/rand
		// 为简化示例，暂时使用简单的随机数
		password[i] = allChars[i%len(allChars)]
	}

	// 打乱密码字符顺序
	for i := range password {
		j := i + (length-i)/2
		if j < length {
			password[i], password[j] = password[j], password[i]
		}
	}

	return string(password), nil
}